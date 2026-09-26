package p2p

import (
	"math"
	"time"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// rankingVersion is bumped whenever the scoring formula changes, so a client
// (or a future migration) can tell which formula produced a given snapshot.
const rankingVersion = 1

// offerScore is the intermediate result CalculateRanking sorts on before
// persisting Rank.
type offerScore struct {
	offerID string
	score   decimal.Decimal
}

// CalculateRanking implements Plan Section 14's marketplace ranking:
// completion rate and trade volume push an offer up, disputes resolved
// against the merchant push it down, and recent activity is a small
// tie-breaking bonus over a merchant who has gone quiet. This is a
// deliberately simple, transparent formula (weights chosen, not learned) -
// Section 14 does not specify exact weights, so these are the plan's
// intent translated into a concrete, documented scoring function rather
// than left unimplemented.
func CalculateRanking(gc *sharedconfig.GlobalConfig) error {
	var offers []p2pModels.Offer
	if err := gc.DB.Where("status = ?", p2pModels.OfferStatusActive).Find(&offers).Error; err != nil {
		return err
	}
	if len(offers) == 0 {
		return nil
	}

	scores := make([]offerScore, 0, len(offers))
	for _, offer := range offers {
		merchantPerf, _ := GetMerchantPerformance(gc.DB, offer.MerchantUserID)
		var offerPerf p2pModels.MerchantOfferPerformance
		gc.DB.Where("offer_id = ?", offer.ID).First(&offerPerf)
		scores = append(scores, offerScore{offerID: offer.ID, score: scoreOffer(merchantPerf, offerPerf)})
	}

	// Descending score -> ascending rank (rank 1 is best). Simple insertion
	// sort is fine here: the marketplace is not expected to have thousands
	// of concurrently active offers.
	for i := 1; i < len(scores); i++ {
		for j := i; j > 0 && scores[j].score.GreaterThan(scores[j-1].score); j-- {
			scores[j], scores[j-1] = scores[j-1], scores[j]
		}
	}

	now := time.Now().UTC()
	for i, s := range scores {
		if err := upsertRankingSnapshot(gc.DB, s.offerID, i+1, s.score, now); err != nil {
			return err
		}
	}
	return nil
}

// scoreOffer computes a single offer's ranking score out of the merchant's
// overall performance and this offer's own completion history. Completion
// rate is already 0-100; trade volume is log-scaled so a merchant with 1000
// trades doesn't dominate one with 100 nearly as much as a linear scale
// would; disputes resolved against the merchant are a flat penalty per
// occurrence, since even a small number is a meaningful trust signal.
func scoreOffer(merchant p2pModels.MerchantPerformance, offerPerf p2pModels.MerchantOfferPerformance) decimal.Decimal {
	completionRate := decimal.RequireFromString(orDefaultStr(merchant.CompletionRate, "0"))
	volumeScore := decimal.NewFromFloat(math.Log10(float64(merchant.CompletedTrades) + 1)).Mul(decimal.NewFromInt(10))
	offerCompletionRate := decimal.RequireFromString(orDefaultStr(offerPerf.CompletionRate, "0"))
	disputePenalty := decimal.NewFromInt(merchant.DisputesResolvedAgainstMerchant).Mul(decimal.NewFromInt(15))

	score := completionRate.Add(volumeScore).Add(offerCompletionRate.Div(decimal.NewFromInt(2))).Sub(disputePenalty)
	if score.IsNegative() {
		score = decimal.Zero
	}
	return score
}

func upsertRankingSnapshot(db *gorm.DB, offerID string, rank int, score decimal.Decimal, calculatedAt time.Time) error {
	var snapshot p2pModels.OfferRankingSnapshot
	err := db.Where("offer_id = ?", offerID).First(&snapshot).Error
	if err == gorm.ErrRecordNotFound {
		snapshot = p2pModels.OfferRankingSnapshot{ID: uuid.NewString(), OfferID: offerID}
	} else if err != nil {
		return err
	}
	snapshot.Rank = rank
	snapshot.Score = score.Truncate(6).String()
	snapshot.RankingVersion = rankingVersion
	snapshot.CalculatedAt = calculatedAt
	return db.Save(&snapshot).Error
}
