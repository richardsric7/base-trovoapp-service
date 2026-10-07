import 'package:flutter_test/flutter_test.dart';
import 'package:trovo_app/network/market_requests.dart';

void main() {
  final pairJson = {
    'base': {'code': 'TROV', 'name': 'Trovo Token', 'contractAddress': '0xAAA', 'imageUrl': ''},
    'counter': {'code': 'CNGN', 'name': 'cNGN', 'contractAddress': '0xBBB', 'imageUrl': ''},
    'lastPrice': '152.5',
    'changePercent24h': '-1.25',
    'high24h': '160',
    'low24h': '150',
    'baseVolume24h': '1200',
    'counterVolume24h': '183000',
    'trades24h': 14,
    'bestAsk': '153',
    'bestBid': '152',
    'lastTradeAt': 1700000000,
  };

  test('pair parses and survives the round trip through viewData', () {
    final p = MarketPair(pairJson);
    expect(p.symbol, 'TROV/CNGN');
    expect(p.isUp, isFalse);
    expect(p.lastTradeAt, DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000));
    final again = MarketPair(p.toMap());
    expect(again.symbol, p.symbol);
    expect(again.changePercent24h, p.changePercent24h);
    expect(again.bestAsk, '153');
    expect(MarketFavorites.id(again), '0xaaa/0xbbb');
  });

  test('order book spread', () {
    final b = MarketOrderBook({
      'asks': [{'price': '153', 'amount': '2'}],
      'bids': [{'price': '152', 'amount': '3'}],
    });
    expect(b.spread, 1);
    expect(b.bids.first.total, 456);
    expect(MarketOrderBook({'asks': [], 'bids': []}).spread, isNull);
  });

  MarketOffer offer(Map<String, dynamic> m) => MarketOffer({
        'id': 'o1',
        'offerType': 'BUY',
        'contractAddress': '0xAAA',
        'currencyIssuer': '0xBBB',
        'pricePerUnit': '150',
        'quantity': '10',
        'remaining': '0',
        'open': false,
        'canceled': 0,
        'sold': '0',
        'received': '0',
        'fills': 0,
        ...m,
      });

  test('offer status and fills', () {
    expect(offer({}).status, 'Pending');
    expect(offer({'open': true, 'remaining': '1500'}).status, 'Open');
    expect(offer({'open': true, 'remaining': '1500'}).cancellable, isTrue);
    // a buy offer sells CNGN and receives TROV
    final part = offer({'open': true, 'remaining': '750', 'sold': '750', 'received': '5', 'fills': 2});
    expect(part.status, 'Part filled');
    expect(part.filledBase, 5);
    expect(part.filledCounter, 750);
    expect(part.averagePrice, 150);
    expect(offer({'blockchainOfferId': '7', 'fills': 1, 'sold': '1500', 'received': '10'}).status, 'Filled');
    expect(offer({'canceled': 1}).status, 'Cancelled');
    expect(offer({'canceled': 1, 'fills': 1}).status, 'Cancelled (part filled)');
    expect(offer({'canceled': 1, 'open': true, 'remaining': '5'}).cancellable, isFalse);
    expect(offer({}).forPair(MarketPair(pairJson)), isTrue);
  });

  test('numbers', () {
    expect(formatMarketNumber(0), '0');
    expect(formatMarketNumber(1234567.5), '1,234,567.50');
    expect(formatMarketNumber(152), '152');
    expect(formatMarketNumber(0.000312345), '0.0003123');
    expect(formatMarketString('not a number'), '0');
  });
}
