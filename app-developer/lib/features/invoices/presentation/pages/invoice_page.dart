import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_dimensions.dart';

/// Halaman pelacakan invoice klien.
/// Menampilkan daftar invoice per perusahaan beserta status pembayaran.
class InvoicePage extends StatefulWidget {
  const InvoicePage({super.key});

  @override
  State<InvoicePage> createState() => _InvoicePageState();
}

class _InvoicePageState extends State<InvoicePage> {
  _InvoiceFilter _activeFilter = _InvoiceFilter.semua;

  static final _currencyFmt =
      NumberFormat.currency(locale: 'id_ID', symbol: 'Rp ', decimalDigits: 0);

  // Data mock — ganti dengan data dari API saat tersedia
  final List<_InvoiceItem> _mockItems = const [
    _InvoiceItem(
      id: 'INV-2026-001',
      companyName: 'PT Maju Bersama',
      amount: 5500000,
      dueDate: '2026-03-30',
      status: _InvoiceStatus.belumDibayar,
      plan: 'Professional',
    ),
    _InvoiceItem(
      id: 'INV-2026-002',
      companyName: 'CV Teknologi Nusantara',
      amount: 2750000,
      dueDate: '2026-03-15',
      status: _InvoiceStatus.jatuhTempo,
      plan: 'Starter',
    ),
    _InvoiceItem(
      id: 'INV-2026-003',
      companyName: 'PT Sinar Harapan',
      amount: 11000000,
      dueDate: '2026-03-20',
      status: _InvoiceStatus.lunas,
      plan: 'Enterprise',
    ),
    _InvoiceItem(
      id: 'INV-2026-004',
      companyName: 'UD Karya Mandiri',
      amount: 2750000,
      dueDate: '2026-04-05',
      status: _InvoiceStatus.belumDibayar,
      plan: 'Starter',
    ),
    _InvoiceItem(
      id: 'INV-2025-089',
      companyName: 'PT Global Solusi',
      amount: 5500000,
      dueDate: '2025-12-31',
      status: _InvoiceStatus.lunas,
      plan: 'Professional',
    ),
  ];

  List<_InvoiceItem> get _filtered {
    if (_activeFilter == _InvoiceFilter.semua) return _mockItems;
    return _mockItems
        .where((i) => i.status == _activeFilter.toStatus())
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    final filtered = _filtered;

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: _buildAppBar(context),
      body: Column(
        children: [
          _buildSummaryBanner(context),
          _buildFilterChips(context),
          Expanded(
            child: filtered.isEmpty
                ? const _EmptyState()
                : _buildList(context, filtered),
          ),
        ],
      ),
    );
  }

  AppBar _buildAppBar(BuildContext context) {
    return AppBar(
      backgroundColor: AppColors.surface,
      title: const Text('Invoice'),
      actions: [
        IconButton(
          icon: const Icon(Icons.search_rounded),
          tooltip: 'Cari invoice',
          onPressed: () => _showComingSoon(context, 'Pencarian invoice'),
        ),
        IconButton(
          icon: const Icon(Icons.filter_list_rounded),
          tooltip: 'Filter lanjutan',
          onPressed: () => _showComingSoon(context, 'Filter lanjutan'),
        ),
      ],
    );
  }

  Widget _buildSummaryBanner(BuildContext context) {
    final total = _mockItems.fold<int>(0, (s, i) => s + i.amount);
    final unpaid = _mockItems
        .where((i) =>
            i.status == _InvoiceStatus.belumDibayar ||
            i.status == _InvoiceStatus.jatuhTempo)
        .fold<int>(0, (s, i) => s + i.amount);
    final overdue = _mockItems
        .where((i) => i.status == _InvoiceStatus.jatuhTempo)
        .length;

    return Container(
      margin: const EdgeInsets.fromLTRB(
        AppDimensions.screenPaddingH,
        AppDimensions.spacing16,
        AppDimensions.screenPaddingH,
        AppDimensions.spacing4,
      ),
      padding: const EdgeInsets.all(AppDimensions.spacing16),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [AppColors.primary800, AppColors.primary600],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(AppDimensions.radiusXl),
      ),
      child: Row(
        children: [
          Expanded(
            child: _SummaryTile(
              label: 'Total Tagihan',
              value: _currencyFmt.format(total),
              valueColor: Colors.white,
            ),
          ),
          Container(
            width: 1,
            height: 36,
            color: Colors.white.withValues(alpha: 0.2),
          ),
          Expanded(
            child: _SummaryTile(
              label: 'Belum Lunas',
              value: _currencyFmt.format(unpaid),
              valueColor: AppColors.secondaryLight,
            ),
          ),
          Container(
            width: 1,
            height: 36,
            color: Colors.white.withValues(alpha: 0.2),
          ),
          Expanded(
            child: _SummaryTile(
              label: 'Jatuh Tempo',
              value: '$overdue invoice',
              valueColor:
                  overdue > 0 ? AppColors.error : AppColors.successLight,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFilterChips(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(
        horizontal: AppDimensions.screenPaddingH,
        vertical: AppDimensions.spacing12,
      ),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          children: _InvoiceFilter.values.map((filter) {
            final isActive = _activeFilter == filter;
            return Padding(
              padding: const EdgeInsets.only(right: AppDimensions.spacing8),
              child: AnimatedContainer(
                duration: const Duration(milliseconds: 200),
                child: GestureDetector(
                  onTap: () => setState(() => _activeFilter = filter),
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: AppDimensions.spacing16,
                      vertical: AppDimensions.spacing8,
                    ),
                    decoration: BoxDecoration(
                      color: isActive
                          ? AppColors.primary700
                          : AppColors.surface,
                      borderRadius:
                          BorderRadius.circular(AppDimensions.radiusFull),
                      border: Border.all(
                        color: isActive
                            ? AppColors.primary700
                            : AppColors.neutral200,
                      ),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        if (filter.icon != null) ...[
                          Icon(
                            filter.icon,
                            size: 14,
                            color: isActive
                                ? Colors.white
                                : filter.iconColor,
                          ),
                          const SizedBox(width: AppDimensions.spacing4),
                        ],
                        Text(
                          filter.label,
                          style: TextStyle(
                            fontSize: 13,
                            fontWeight: isActive
                                ? FontWeight.w600
                                : FontWeight.w400,
                            color: isActive
                                ? Colors.white
                                : AppColors.neutral700,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            );
          }).toList(),
        ),
      ),
    );
  }

  Widget _buildList(BuildContext context, List<_InvoiceItem> items) {
    return ListView.separated(
      padding: const EdgeInsets.fromLTRB(
        AppDimensions.screenPaddingH,
        0,
        AppDimensions.screenPaddingH,
        AppDimensions.spacing80,
      ),
      itemCount: items.length,
      separatorBuilder: (_, __) =>
          const SizedBox(height: AppDimensions.spacing12),
      itemBuilder: (context, index) => _InvoiceCard(
        item: items[index],
        currencyFmt: _currencyFmt,
      ),
    );
  }

  void _showComingSoon(BuildContext context, String feature) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Fitur $feature belum tersedia')),
    );
  }
}

// ── Summary Tile ──────────────────────────────────────────────────────────────

class _SummaryTile extends StatelessWidget {
  final String label;
  final String value;
  final Color valueColor;

  const _SummaryTile({
    required this.label,
    required this.value,
    required this.valueColor,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Text(
          label,
          style: TextStyle(
            color: Colors.white.withValues(alpha: 0.65),
            fontSize: 10,
            fontWeight: FontWeight.w500,
            letterSpacing: 0.3,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: AppDimensions.spacing4),
        Text(
          value,
          style: TextStyle(
            color: valueColor,
            fontSize: 12,
            fontWeight: FontWeight.w700,
          ),
          textAlign: TextAlign.center,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
        ),
      ],
    );
  }
}

// ── Invoice Card ──────────────────────────────────────────────────────────────

class _InvoiceCard extends StatelessWidget {
  final _InvoiceItem item;
  final NumberFormat currencyFmt;

  const _InvoiceCard({required this.item, required this.currencyFmt});

  @override
  Widget build(BuildContext context) {
    final (statusColor, statusBg, statusLabel, statusIcon) =
        item.status.display;

    return Container(
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(AppDimensions.radiusXl),
        border: Border.all(
          color: item.status == _InvoiceStatus.jatuhTempo
              ? AppColors.errorBase.withValues(alpha: 0.3)
              : AppColors.neutral200,
        ),
        boxShadow: const [
          BoxShadow(
            color: Color(0x06000000),
            blurRadius: 8,
            offset: Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        children: [
          // ── Header ──
          Padding(
            padding: const EdgeInsets.all(AppDimensions.cardPadding),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Icon plan
                Container(
                  width: 44,
                  height: 44,
                  decoration: BoxDecoration(
                    color: AppColors.primary50,
                    borderRadius:
                        BorderRadius.circular(AppDimensions.radiusMd),
                  ),
                  child: const Icon(
                    Icons.receipt_long_rounded,
                    color: AppColors.primary700,
                    size: AppDimensions.iconMd,
                  ),
                ),
                const SizedBox(width: AppDimensions.spacing12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        item.companyName,
                        style: Theme.of(context)
                            .textTheme
                            .titleSmall
                            ?.copyWith(color: AppColors.neutral900),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: AppDimensions.spacing2),
                      Row(
                        children: [
                          Text(
                            item.id,
                            style: const TextStyle(
                              fontSize: 12,
                              color: AppColors.neutral500,
                              fontFamily: 'JetBrainsMono',
                            ),
                          ),
                          const SizedBox(width: AppDimensions.spacing8),
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: AppDimensions.spacing6,
                              vertical: 1,
                            ),
                            decoration: BoxDecoration(
                              color: AppColors.accent50,
                              borderRadius: BorderRadius.circular(
                                  AppDimensions.radiusFull),
                            ),
                            child: Text(
                              item.plan,
                              style: const TextStyle(
                                fontSize: 10,
                                color: AppColors.accent400,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                // Status badge
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: AppDimensions.spacing8,
                    vertical: AppDimensions.spacing4,
                  ),
                  decoration: BoxDecoration(
                    color: statusBg,
                    borderRadius:
                        BorderRadius.circular(AppDimensions.radiusFull),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(statusIcon, size: 12, color: statusColor),
                      const SizedBox(width: 3),
                      Text(
                        statusLabel,
                        style: TextStyle(
                          fontSize: 11,
                          color: statusColor,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          // ── Divider ──
          const Divider(height: 1, color: AppColors.neutral100),
          // ── Footer: nominal + due date ──
          Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: AppDimensions.cardPadding,
              vertical: AppDimensions.spacing12,
            ),
            child: Row(
              children: [
                const Icon(Icons.payments_outlined,
                    size: 16, color: AppColors.neutral400),
                const SizedBox(width: AppDimensions.spacing6),
                Text(
                  currencyFmt.format(item.amount),
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        color: AppColors.primary700,
                        fontWeight: FontWeight.w700,
                      ),
                ),
                const Spacer(),
                const Icon(Icons.calendar_today_outlined,
                    size: 14, color: AppColors.neutral400),
                const SizedBox(width: AppDimensions.spacing4),
                Text(
                  'Jatuh tempo: ${_formatDate(item.dueDate)}',
                  style: TextStyle(
                    fontSize: 12,
                    color: item.status == _InvoiceStatus.jatuhTempo
                        ? AppColors.errorBase
                        : AppColors.neutral500,
                    fontWeight: item.status == _InvoiceStatus.jatuhTempo
                        ? FontWeight.w600
                        : FontWeight.w400,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  String _formatDate(String iso) {
    final dt = DateTime.tryParse(iso);
    if (dt == null) return iso;
    return DateFormat('d MMM yyyy', 'id_ID').format(dt);
  }
}

// ── Empty State ───────────────────────────────────────────────────────────────

class _EmptyState extends StatelessWidget {
  const _EmptyState();

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 80,
            height: 80,
            decoration: const BoxDecoration(
              color: AppColors.primary50,
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.receipt_long_outlined,
              size: 40,
              color: AppColors.primary600,
            ),
          ),
          const SizedBox(height: AppDimensions.spacing16),
          Text(
            'Tidak Ada Invoice',
            style: Theme.of(context)
                .textTheme
                .titleMedium
                ?.copyWith(color: AppColors.neutral900),
          ),
          const SizedBox(height: AppDimensions.spacing8),
          const Text(
            'Tidak ada invoice dengan filter ini.',
            style: TextStyle(color: AppColors.neutral500, fontSize: 14),
          ),
        ],
      ),
    );
  }
}

// ── Models ────────────────────────────────────────────────────────────────────

enum _InvoiceStatus {
  belumDibayar,
  jatuhTempo,
  lunas;

  (Color, Color, String, IconData) get display => switch (this) {
        _InvoiceStatus.belumDibayar => (
            AppColors.warningBase,
            AppColors.warningLight,
            'Belum Dibayar',
            Icons.schedule_rounded,
          ),
        _InvoiceStatus.jatuhTempo => (
            AppColors.errorBase,
            AppColors.errorLight,
            'Jatuh Tempo',
            Icons.warning_amber_rounded,
          ),
        _InvoiceStatus.lunas => (
            AppColors.successBase,
            AppColors.successLight,
            'Lunas',
            Icons.check_circle_outline_rounded,
          ),
      };
}

enum _InvoiceFilter {
  semua,
  belumDibayar,
  jatuhTempo,
  lunas;

  String get label => switch (this) {
        _InvoiceFilter.semua => 'Semua',
        _InvoiceFilter.belumDibayar => 'Belum Dibayar',
        _InvoiceFilter.jatuhTempo => 'Jatuh Tempo',
        _InvoiceFilter.lunas => 'Lunas',
      };

  IconData? get icon => switch (this) {
        _InvoiceFilter.semua => null,
        _InvoiceFilter.belumDibayar => Icons.schedule_rounded,
        _InvoiceFilter.jatuhTempo => Icons.warning_amber_rounded,
        _InvoiceFilter.lunas => Icons.check_circle_outline_rounded,
      };

  Color get iconColor => switch (this) {
        _InvoiceFilter.semua => AppColors.neutral500,
        _InvoiceFilter.belumDibayar => AppColors.warningBase,
        _InvoiceFilter.jatuhTempo => AppColors.errorBase,
        _InvoiceFilter.lunas => AppColors.successBase,
      };

  _InvoiceStatus? toStatus() => switch (this) {
        _InvoiceFilter.semua => null,
        _InvoiceFilter.belumDibayar => _InvoiceStatus.belumDibayar,
        _InvoiceFilter.jatuhTempo => _InvoiceStatus.jatuhTempo,
        _InvoiceFilter.lunas => _InvoiceStatus.lunas,
      };
}

class _InvoiceItem {
  final String id;
  final String companyName;
  final int amount;
  final String dueDate;
  final _InvoiceStatus status;
  final String plan;

  const _InvoiceItem({
    required this.id,
    required this.companyName,
    required this.amount,
    required this.dueDate,
    required this.status,
    required this.plan,
  });
}
