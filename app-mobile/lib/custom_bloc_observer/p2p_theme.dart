import 'package:flutter/material.dart';

// P2PTheme is a small, real design-token set scoped to the new P2P screens
// only (Plan Section 94). Neither the rest of this app nor app-web has a
// consistently-enforced design system today, so rather than perpetuate the
// existing ad hoc `height/N` spacing division or invent unrelated new
// colors, this formalizes values already in use somewhere in the product:
// p2pBrandDark matches this app's own `trovoblue` (0xFF00225A), which is
// also the undocumented second brand blue app-web's own best-in-class page
// independently uses - the two platforms should read as one product.
class P2PTheme {
  static const Color brandDark = Color(0xFF00225A);
  static const Color primary = Color(0xFF004988);
  static const Color success = Color(0xFF00A859);
  static const Color danger = Color(0xFFBE3800);
  static const Color warning = Color(0xFFB7791F);
  static const Color neutralBg = Color(0xFFF2F6F9);

  // Fixed 4pt spacing scale, replacing this feature's use of the app's
  // legacy `height/N` ad hoc division.
  static const double space1 = 4;
  static const double space2 = 8;
  static const double space3 = 12;
  static const double space4 = 16;
  static const double space6 = 24;
  static const double space8 = 32;

  // A deliberate middle value between app-web's 32px best-in-class panel
  // radius and this app's own dominant 15px card radius - reads as
  // family resemblance to both without copying either exactly.
  static const double cardRadius = 16;
  static const double sheetRadius = 24;

  // Matches app-web's history page shadow value exactly, so both
  // platforms' P2P cards read as the same product.
  static List<BoxShadow> cardShadow = [
    BoxShadow(
      color: const Color(0x0F00225A),
      blurRadius: 60,
      offset: const Offset(0, 18),
    ),
  ];

  static Color statusColor(String orderStatus) {
    switch (orderStatus) {
      case 'AWAITING_ESCROW_DEPOSIT':
      case 'AWAITING_PAYMENT':
        return warning;
      case 'AWAITING_PAYMENT_CONFIRMATION':
        return primary;
      case 'COMPLETED':
        return success;
      case 'REJECTED':
      case 'CANCELLED':
      case 'EXPIRED':
        return const Color(0xFF9EA3AE); // matches existing `grey`
      default:
        return const Color(0xFF9EA3AE);
    }
  }

  static String statusLabel(String orderStatus) {
    switch (orderStatus) {
      case 'AWAITING_APPROVAL':
        return 'Awaiting approval';
      case 'AWAITING_ESCROW_DEPOSIT':
        return 'Awaiting escrow deposit';
      case 'AWAITING_PAYMENT':
        return 'Awaiting payment';
      case 'AWAITING_PAYMENT_CONFIRMATION':
        return 'Awaiting confirmation';
      case 'COMPLETED':
        return 'Completed';
      case 'REJECTED':
        return 'Rejected';
      case 'CANCELLED':
        return 'Cancelled';
      case 'EXPIRED':
        return 'Expired';
      default:
        return orderStatus;
    }
  }
}

// P2PStatusPill is the shared status-pill component (Plan Section 94) -
// used consistently everywhere an order/dispute status is shown.
class P2PStatusPill extends StatelessWidget {
  final String status;
  final bool disputed;
  const P2PStatusPill({super.key, required this.status, this.disputed = false});

  @override
  Widget build(BuildContext context) {
    final color = disputed ? P2PTheme.danger : P2PTheme.statusColor(status);
    final label = disputed ? 'Disputed' : P2PTheme.statusLabel(status);
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: P2PTheme.space3,
        vertical: P2PTheme.space1,
      ),
      decoration: BoxDecoration(
        color: color.withOpacity(0.12),
        borderRadius: BorderRadius.circular(P2PTheme.space6),
      ),
      child: Text(
        label,
        style: TextStyle(
          color: color,
          fontSize: 12,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }
}

// P2PEmptyState is the shared empty-state component (Plan Section 97.1) -
// an icon + message + primary CTA, replacing this app's plain
// centered-text-only empty states.
class P2PEmptyState extends StatelessWidget {
  final IconData icon;
  final String message;
  final String? ctaLabel;
  final VoidCallback? onCta;
  const P2PEmptyState({
    super.key,
    required this.icon,
    required this.message,
    this.ctaLabel,
    this.onCta,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(P2PTheme.space8),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(icon, size: 48, color: P2PTheme.primary.withOpacity(0.4)),
            const SizedBox(height: P2PTheme.space4),
            Text(
              message,
              textAlign: TextAlign.center,
              style: const TextStyle(color: Colors.black54, fontSize: 15),
            ),
            if (ctaLabel != null) ...[
              const SizedBox(height: P2PTheme.space4),
              TextButton(onPressed: onCta, child: Text(ctaLabel!)),
            ],
          ],
        ),
      ),
    );
  }
}

// P2PListCard is the shared offer/order list-tile card (Plan Section 94).
class P2PListCard extends StatelessWidget {
  final Widget child;
  final VoidCallback? onTap;
  const P2PListCard({super.key, required this.child, this.onTap});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: P2PTheme.space3),
      child: Material(
        color: Colors.white,
        borderRadius: BorderRadius.circular(P2PTheme.cardRadius),
        child: InkWell(
          borderRadius: BorderRadius.circular(P2PTheme.cardRadius),
          onTap: onTap,
          child: Container(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(P2PTheme.cardRadius),
              boxShadow: P2PTheme.cardShadow,
            ),
            padding: const EdgeInsets.all(P2PTheme.space4),
            child: child,
          ),
        ),
      ),
    );
  }
}
