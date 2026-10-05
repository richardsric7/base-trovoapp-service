import 'package:easy_localization/easy_localization.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart' as intl;
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/public_market.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';

// Shared pieces of the Public Markets screens (Trovo_App_Stocks prototype).

final _money = intl.NumberFormat('#,##0.00');
final _qty = intl.NumberFormat('#,##0.####');

String pmMoney(dynamic v, [String code = 'NGN']) => v == null || '$v'.isEmpty ? '—' : '${_money.format(pmNum(v))} $code';
String pmNaira(dynamic v) => v == null || '$v'.isEmpty ? '—' : '₦${_money.format(pmNum(v))}';
String pmQty(dynamic v) => v == null || '$v'.isEmpty ? '—' : _qty.format(pmNum(v));
String pmDate(String iso) {
  final t = DateTime.tryParse(iso);
  if (t == null || t.year < 2000) return '—';
  return intl.DateFormat('d MMM yyyy').format(t.toLocal());
}

String pmDateTime(String iso) {
  final t = DateTime.tryParse(iso);
  if (t == null || t.year < 2000) return '—';
  return intl.DateFormat('d MMM yyyy, HH:mm').format(t.toLocal());
}

// pmOpen pushes a Public Markets page with its arguments.
void pmOpen(DataProvider appState, PageConfiguration page, [Map? args]) {
  appState.viewData ??= {};
  if (args != null) appState.viewData![page.key] = args;
  appState.currentAction = PageAction(state: PageState.addPage, page: page);
}

Map pmArgs(DataProvider appState, PageConfiguration page) => (appState.viewData?[page.key] as Map?) ?? {};

class PMLogo extends StatelessWidget {
  final PMAsset asset;
  final double size;
  const PMLogo({super.key, required this.asset, this.size = 44});

  Color _hex(String v, Color fallback) {
    final h = v.replaceAll('#', '');
    final n = int.tryParse(h.length == 6 ? 'FF$h' : h, radix: 16);
    return n == null ? fallback : Color(n);
  }

  @override
  Widget build(BuildContext context) {
    final bg = _hex('${asset.logo['background'] ?? ''}', P2PTheme.brandDark);
    final fg = _hex('${asset.logo['foreground'] ?? ''}', Colors.white);
    final url = '${asset.logo['url'] ?? ''}';
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(color: bg, borderRadius: BorderRadius.circular(size / 3.2)),
      clipBehavior: Clip.antiAlias,
      alignment: Alignment.center,
      child: url.startsWith('http')
          ? Image.network(url, fit: BoxFit.cover, errorBuilder: (_, __, ___) => _text(fg))
          : _text(fg),
    );
  }

  Widget _text(Color fg) => Text(asset.initials, style: TextStyle(color: fg, fontWeight: FontWeight.w700, fontSize: size / 3));
}

class PMChange extends StatelessWidget {
  final String percent;
  final double fontSize;
  const PMChange({super.key, required this.percent, this.fontSize = 13});

  @override
  Widget build(BuildContext context) {
    final v = pmNum(percent);
    final color = v > 0 ? P2PTheme.success : (v < 0 ? P2PTheme.danger : Colors.black45);
    final arrow = v > 0 ? '▲' : (v < 0 ? '▼' : '');
    return Text('$arrow ${v.abs().toStringAsFixed(2)}%', style: TextStyle(color: color, fontSize: fontSize, fontWeight: FontWeight.w600));
  }
}

class PMBadge extends StatelessWidget {
  final String label;
  final Color color;
  const PMBadge(this.label, {super.key, this.color = P2PTheme.primary});

  @override
  Widget build(BuildContext context) => Container(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
        decoration: BoxDecoration(color: color.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(20)),
        child: Text(label, style: TextStyle(color: color, fontSize: 11, fontWeight: FontWeight.w600)),
      );
}

class PMSection extends StatelessWidget {
  final String title;
  final Widget child;
  final Widget? trailing;
  const PMSection({super.key, required this.title, required this.child, this.trailing});

  @override
  Widget build(BuildContext context) => Container(
        margin: const EdgeInsets.only(bottom: P2PTheme.space4),
        padding: const EdgeInsets.all(P2PTheme.space4),
        decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(P2PTheme.cardRadius), boxShadow: P2PTheme.cardShadow),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Row(children: [
            Expanded(child: Text(title, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 16, color: P2PTheme.brandDark))),
            if (trailing != null) trailing!,
          ]),
          const SizedBox(height: P2PTheme.space3),
          child,
        ]),
      );
}

class PMRow extends StatelessWidget {
  final String label;
  final String value;
  final Color? valueColor;
  const PMRow(this.label, this.value, {super.key, this.valueColor});

  @override
  Widget build(BuildContext context) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 6),
        child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Expanded(child: Text(label, style: const TextStyle(color: Colors.black54, fontSize: 13))),
          const SizedBox(width: 12),
          Flexible(
            child: Text(value.isEmpty ? '—' : value,
                textAlign: TextAlign.right, style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: valueColor ?? P2PTheme.brandDark)),
          ),
        ]),
      );
}

// PMChart is the price line of a range.
class PMChart extends StatelessWidget {
  final List<Map> points;
  final double height;
  const PMChart({super.key, required this.points, this.height = 180});

  @override
  Widget build(BuildContext context) {
    final values = points.map((p) => pmNum(p['price'])).toList();
    if (values.length < 2) {
      return SizedBox(height: height, child: Center(child: Text('pmnochart'.tr(), style: const TextStyle(color: Colors.black45))));
    }
    final up = values.last >= values.first;
    final color = up ? P2PTheme.success : P2PTheme.danger;
    final minY = values.reduce((a, b) => a < b ? a : b);
    final maxY = values.reduce((a, b) => a > b ? a : b);
    final pad = (maxY - minY).abs() < 0.0001 ? 1.0 : (maxY - minY) * 0.1;
    return SizedBox(
      height: height,
      child: LineChart(LineChartData(
        minY: minY - pad,
        maxY: maxY + pad,
        gridData: const FlGridData(show: false),
        titlesData: const FlTitlesData(show: false),
        borderData: FlBorderData(show: false),
        lineTouchData: LineTouchData(
          touchTooltipData: LineTouchTooltipData(
            getTooltipItems: (spots) => spots
                .map((s) => LineTooltipItem(
                    '${pmNaira(values[s.x.toInt()])}\n${pmDateTime('${points[s.x.toInt()]['at']}')}', const TextStyle(color: Colors.white, fontSize: 11)))
                .toList(),
          ),
        ),
        lineBarsData: [
          LineChartBarData(
            spots: [for (var i = 0; i < values.length; i++) FlSpot(i.toDouble(), values[i])],
            isCurved: true,
            color: color,
            barWidth: 2,
            dotData: const FlDotData(show: false),
            belowBarData: BarAreaData(show: true, color: color.withValues(alpha: 0.08)),
          ),
        ],
      )),
    );
  }
}

// pmStateLabel is an order state in words.
String pmStateLabel(String state) {
  switch (state) {
    case 'awaiting-payment':
      return 'pmstateawaitingpayment'.tr();
    case 'queued':
    case 'filled-from-inventory':
      return 'pmstatequeued'.tr();
    case 'pending-execution':
      return 'pmstatependingexecution'.tr();
    case 'executed':
      return 'pmstateexecuted'.tr();
    case 'settlement_final':
    case 'submitted':
    case 'chain_final':
      return 'pmstatesettling'.tr();
    case 'complete':
      return 'pmstatecomplete'.tr();
    case 'rejected':
      return 'pmstaterejected'.tr();
    case 'failed':
      return 'pmstatefailed'.tr();
    case 'cancelled':
      return 'pmstatecancelled'.tr();
  }
  return state;
}

Color pmStateColor(String state) {
  if (state == 'complete') return P2PTheme.success;
  if (const ['rejected', 'failed', 'cancelled'].contains(state)) return P2PTheme.danger;
  return P2PTheme.warning;
}
