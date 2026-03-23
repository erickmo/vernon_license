import 'package:flutter/material.dart';
import '../../../../../core/constants/app_dimensions.dart';

class QuickActionWidget extends StatelessWidget {
  final IconData icon;
  final Color iconColor;
  final Color iconBackground;
  final String label;
  final VoidCallback onTap;

  const QuickActionWidget({
    super.key,
    required this.icon,
    required this.iconColor,
    required this.iconBackground,
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          _buildIconContainer(),
          const SizedBox(height: AppDimensions.spacing8),
          _buildLabel(context),
        ],
      ),
    );
  }

  Widget _buildIconContainer() => Container(
        width: 56,
        height: 56,
        decoration: BoxDecoration(
          color: iconBackground,
          borderRadius: BorderRadius.circular(AppDimensions.radiusXl),
        ),
        child: Icon(icon, color: iconColor, size: AppDimensions.iconLg),
      );

  Widget _buildLabel(BuildContext context) => SizedBox(
        width: 72,
        child: Text(
          label,
          style: Theme.of(context).textTheme.labelMedium,
          textAlign: TextAlign.center,
          maxLines: 2,
          overflow: TextOverflow.ellipsis,
        ),
      );
}
