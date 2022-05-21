import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';

LineChartData mainData(color,opcity) {


  List<Color> gradientColors = [color];
  return LineChartData(
    gridData: FlGridData(
      show: false,
      drawVerticalLine: true,
      getDrawingHorizontalLine: (value) {
        return FlLine(strokeWidth: 1, color: Colors.transparent);
      },
      getDrawingVerticalLine: (value) {
        return FlLine(
          strokeWidth: 0,
        );
      },
    ),
    titlesData: FlTitlesData(
      show: false,
      leftTitles: SideTitles(
        showTitles: false,
        reservedSize: 20,
        margin: 8,
      ),
    ),
    borderData: FlBorderData(
        show: true,
        border: Border.all(color: Colors.white, width: 0)),
    minX: 0,
    maxX: 8,
    minY: 0,
    maxY: 5,
    lineBarsData: [
      LineChartBarData(
        spots: [
          const FlSpot(0, 2.5),
          const FlSpot(1, 2),
          const FlSpot(2, 4),
          const FlSpot(3, 3.1),
          const FlSpot(4, 4),
          const FlSpot(5, 2),
          const FlSpot(6, 4),
          const FlSpot(7, 3.1),
          const FlSpot(8, 2),
          const FlSpot(9, 1.5),
          const FlSpot(10, 3),
        ],
        isCurved: true,
        colors: gradientColors,
        barWidth: 2,
        isStrokeCapRound: false,
        dotData: FlDotData(
          show: false,
        ),
        belowBarData: BarAreaData(
          show: true,
          colors:
          gradientColors.map((color) => color.withOpacity(opcity)).toList(),
        ),
      ),
    ],
  );
}