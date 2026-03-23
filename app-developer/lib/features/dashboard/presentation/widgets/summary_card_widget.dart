import 'package:flutter/material.dart';
import '../../../../../core/constants/app_colors.dart';
import '../../../../../core/constants/app_dimensions.dart';

class SummaryCardWidget extends StatelessWidget {
  final IconData icon;
  final Color iconColor;
  final Color iconBackground;
  final String label;
  final String value;
  final String? subtitle;
  final VoidCallback? onTap;

  const SummaryCardWidget({
    super.key,
    required this.icon,
    required this.iconColor,
    required this.iconBackground,
    required this.label,
    required this.value,
    this.subtitle,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(AppDimensions.cardPadding),
        decoration: BoxDecoration(
          color: AppColors.surface,
          borderRadius: BorderRadius.circular(AppDimensions.radiusLg),
          border: Border.all(color: AppColors.neutral200),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildIconBox(),
            const SizedBox(height: AppDimensions.spacing12),
            _buildValue(context),
            const SizedBox(height: AppDimensions.spacing4),
            _buildLabel(context),
            if (subtitle != null) ...[
              const SizedBox(height: AppDimensions.spacing4),
              _buildSubtitle(context),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildIconBox() => Container(
        width: 40,
        height: 40,
        decoration: BoxDecoration(
          color: iconBackground,
          borderRadius: BorderRadius.circular(AppDimensions.radiusMd),
        ),
        child: Icon(icon, color: iconColor, size: AppDimensions.iconMd),
      );

  Widget _buildValue(BuildContext context) => Text(
        value,
        style: Theme.of(context)
            .textTheme
            .titleLarge
            ?.copyWith(fontWeight: FontWeight.w700),
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
      );

  Widget _buildLabel(BuildContext context) => Text(
        label,
        style: Theme.of(context).textTheme.bodySmall,
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
      );

  Widget _buildSubtitle(BuildContext context) => Text(
        subtitle!,
        style: Theme.of(context).textTheme.labelSmall,
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
      );
}
