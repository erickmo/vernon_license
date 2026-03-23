import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import '../../../../core/auth/auth_notifier.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_dimensions.dart';
import '../../../../injection_container.dart' show sl;
import '../../domain/entities/developer_dashboard_entity.dart';
import '../cubit/developer_dashboard_cubit.dart';
import '../widgets/quick_action_widget.dart';
import '../widgets/summary_card_widget.dart';

class DashboardPage extends StatelessWidget {
  const DashboardPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<DeveloperDashboardCubit>()..loadDashboard(),
      child: const _DashboardView(),
    );
  }
}

class _DashboardView extends StatelessWidget {
  const _DashboardView();

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<DeveloperDashboardCubit, DeveloperDashboardState>(
      builder: (context, state) {
        if (state is DeveloperDashboardLoading) return const _LoadingView();
        if (state is DeveloperDashboardError) {
          return _ErrorView(message: state.message);
        }
        if (state is DeveloperDashboardLoaded) {
          return _ContentView(data: state.data);
        }
        return const _LoadingView();
      },
    );
  }
}

class _LoadingView extends StatelessWidget {
  const _LoadingView();

  @override
  Widget build(BuildContext context) => const Scaffold(
        body: Center(child: CircularProgressIndicator()),
      );
}

class _ErrorView extends StatelessWidget {
  final String message;
  const _ErrorView({required this.message});

  @override
  Widget build(BuildContext context) => Scaffold(
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline,
                  size: 48, color: AppColors.errorBase),
              const SizedBox(height: AppDimensions.spacing16),
              Text(message, textAlign: TextAlign.center),
              const SizedBox(height: AppDimensions.spacing16),
              FilledButton(
                onPressed: () =>
                    context.read<DeveloperDashboardCubit>().loadDashboard(),
                child: const Text('Coba Lagi'),
              ),
            ],
          ),
        ),
      );
}

class _ContentView extends StatelessWidget {
  final DeveloperDashboardEntity data;

  const _ContentView({required this.data});

  static final _currencyFmt =
      NumberFormat.currency(locale: 'id_ID', symbol: 'Rp ', decimalDigits: 0);

  @override
  Widget build(BuildContext context) {
    final authNotifier = sl<AuthNotifier>();

    return Scaffold(
      body: RefreshIndicator(
        onRefresh: () =>
            context.read<DeveloperDashboardCubit>().loadDashboard(),
        child: CustomScrollView(
          slivers: [
            _buildSliverAppBar(context, authNotifier),
            SliverPadding(
              padding: const EdgeInsets.symmetric(
                horizontal: AppDimensions.screenPaddingH,
                vertical: AppDimensions.screenPaddingV,
              ),
              sliver: SliverList(
                delegate: SliverChildListDelegate([
                  _buildMrrCard(context),
                  const SizedBox(height: AppDimensions.spacing20),
                  _buildKpiRow(context),
                  const SizedBox(height: AppDimensions.spacing20),
                  _buildKpiRow2(context),
                  const SizedBox(height: AppDimensions.spacing24),
                  _buildQuickActions(context),
                  if (data.recentRegistrations.isNotEmpty) ...[
                    const SizedBox(height: AppDimensions.spacing24),
                    _buildRecentRegistrations(context),
                  ],
                  const SizedBox(height: AppDimensions.spacing80),
                ]),
              ),
            ),
          ],
        ),
      ),
    );
  }

  SliverAppBar _buildSliverAppBar(
      BuildContext context, AuthNotifier authNotifier) {
    return SliverAppBar(
      expandedHeight: 180,
      pinned: true,
      backgroundColor: AppColors.primary900,
      foregroundColor: Colors.white,
      elevation: 0,
      title: Row(
        children: [
          Container(
            width: 28,
            height: 28,
            decoration: BoxDecoration(
              color: Colors.white.withValues(alpha: 0.15),
              borderRadius: BorderRadius.circular(6),
            ),
            child: const Icon(Icons.bolt_rounded, color: Colors.white, size: 18),
          ),
          const SizedBox(width: 8),
          const Text(
            'FlashERP',
            style: TextStyle(
              color: Colors.white,
              fontWeight: FontWeight.w700,
              fontSize: 18,
            ),
          ),
        ],
      ),
      actions: [
        IconButton(
          icon: const Icon(Icons.logout_outlined, color: Colors.white),
          tooltip: 'Keluar',
          onPressed: () async => authNotifier.onLogout(),
        ),
      ],
      flexibleSpace: FlexibleSpaceBar(
        collapseMode: CollapseMode.parallax,
        background: _buildAppBarBackground(context, authNotifier),
      ),
    );
  }

  Widget _buildAppBarBackground(
      BuildContext context, AuthNotifier authNotifier) {
    final hour = DateTime.now().hour;
    final greeting = hour < 12
        ? 'Selamat Pagi'
        : hour < 15
            ? 'Selamat Siang'
            : hour < 18
                ? 'Selamat Sore'
                : 'Selamat Malam';
    final roleLabel = _labelForRole(authNotifier.userRole ?? '');
    final initials = _getInitials(roleLabel);

    return Stack(
      fit: StackFit.expand,
      children: [
        // Gradient latar
        Container(
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              colors: [AppColors.primary900, AppColors.primary600],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
        ),
        // Lingkaran dekoratif kanan-atas
        Positioned(
          right: -40,
          top: -20,
          child: Container(
            width: 140,
            height: 140,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: Colors.white.withValues(alpha: 0.05),
            ),
          ),
        ),
        // Lingkaran dekoratif kanan-bawah
        Positioned(
          right: 30,
          bottom: -25,
          child: Container(
            width: 90,
            height: 90,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: AppColors.secondary.withValues(alpha: 0.18),
            ),
          ),
        ),
        // Garis aksen kiri
        Positioned(
          left: 0,
          top: 60,
          child: Container(
            width: 3,
            height: 60,
            decoration: BoxDecoration(
              gradient: LinearGradient(
                colors: [
                  AppColors.secondary.withValues(alpha: 0),
                  AppColors.secondary.withValues(alpha: 0.8),
                  AppColors.secondary.withValues(alpha: 0),
                ],
                begin: Alignment.topCenter,
                end: Alignment.bottomCenter,
              ),
            ),
          ),
        ),
        // Konten — greeting + avatar
        Positioned(
          left: AppDimensions.screenPaddingH,
          right: AppDimensions.screenPaddingH,
          bottom: 20,
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    // Badge role
                    Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: AppDimensions.spacing8,
                          vertical: 3),
                      decoration: BoxDecoration(
                        color: AppColors.secondary.withValues(alpha: 0.18),
                        borderRadius:
                            BorderRadius.circular(AppDimensions.radiusFull),
                        border: Border.all(
                            color: AppColors.secondary.withValues(alpha: 0.4)),
                      ),
                      child: Text(
                        roleLabel,
                        style: const TextStyle(
                          color: AppColors.secondaryLight,
                          fontSize: 11,
                          fontWeight: FontWeight.w600,
                          letterSpacing: 0.6,
                        ),
                      ),
                    ),
                    const SizedBox(height: AppDimensions.spacing8),
                    Text(
                      greeting,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 22,
                        fontWeight: FontWeight.w700,
                        letterSpacing: -0.4,
                        height: 1.1,
                      ),
                    ),
                    const SizedBox(height: AppDimensions.spacing4),
                    Text(
                      DateFormat('EEEE, d MMMM yyyy', 'id_ID')
                          .format(DateTime.now()),
                      style: TextStyle(
                        color: Colors.white.withValues(alpha: 0.6),
                        fontSize: 12,
                        letterSpacing: 0.1,
                      ),
                    ),
                  ],
                ),
              ),
              // Avatar initials
              Container(
                width: 52,
                height: 52,
                decoration: BoxDecoration(
                  gradient: const LinearGradient(
                    colors: [AppColors.secondaryLight, AppColors.secondary400],
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                  ),
                  shape: BoxShape.circle,
                  border: Border.all(
                      color: Colors.white.withValues(alpha: 0.3), width: 2),
                  boxShadow: [
                    BoxShadow(
                      color: AppColors.secondary.withValues(alpha: 0.35),
                      blurRadius: 12,
                      offset: const Offset(0, 4),
                    ),
                  ],
                ),
                child: Center(
                  child: Text(
                    initials,
                    style: const TextStyle(
                      color: AppColors.primary900,
                      fontWeight: FontWeight.w800,
                      fontSize: 17,
                      letterSpacing: 0.5,
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  String _getInitials(String roleLabel) {
    final parts = roleLabel.trim().split(' ');
    if (parts.length >= 2) {
      return '${parts[0][0]}${parts[1][0]}'.toUpperCase();
    }
    return roleLabel.isNotEmpty ? roleLabel[0].toUpperCase() : 'D';
  }

  Widget _buildMrrCard(BuildContext context) {
    final isPositive = data.mrrGrowthPercent >= 0;
    return Container(
      padding: const EdgeInsets.all(AppDimensions.spacing20),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [AppColors.primary700, AppColors.primary500],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(AppDimensions.radiusXl),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Monthly Recurring Revenue',
            style: Theme.of(context)
                .textTheme
                .labelLarge
                ?.copyWith(color: Colors.white70),
          ),
          const SizedBox(height: AppDimensions.spacing8),
          Text(
            _currencyFmt.format(data.mrr),
            style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                  color: Colors.white,
                  fontWeight: FontWeight.w700,
                ),
          ),
          const SizedBox(height: AppDimensions.spacing8),
          Row(
            children: [
              Icon(
                isPositive ? Icons.trending_up : Icons.trending_down,
                color: isPositive
                    ? AppColors.successLight
                    : AppColors.errorLight,
                size: AppDimensions.iconSm,
              ),
              const SizedBox(width: AppDimensions.spacing4),
              Text(
                '${isPositive ? '+' : ''}${data.mrrGrowthPercent.toStringAsFixed(1)}% vs bulan lalu',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: isPositive
                          ? AppColors.successLight
                          : AppColors.errorLight,
                    ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildKpiRow(BuildContext context) => Row(
        children: [
          Expanded(
            child: SummaryCardWidget(
              icon: Icons.business_outlined,
              iconColor: AppColors.accent400,
              iconBackground: AppColors.accent50,
              label: 'Total Klien',
              value: data.totalClients.toString(),
            ),
          ),
          const SizedBox(width: AppDimensions.spacing12),
          Expanded(
            child: SummaryCardWidget(
              icon: Icons.domain_outlined,
              iconColor: AppColors.infoBase,
              iconBackground: AppColors.infoLight,
              label: 'Perusahaan Aktif',
              value: data.activeCompanies.toString(),
            ),
          ),
          const SizedBox(width: AppDimensions.spacing12),
          Expanded(
            child: SummaryCardWidget(
              icon: Icons.receipt_long_outlined,
              iconColor: AppColors.successBase,
              iconBackground: AppColors.successLight,
              label: 'Faktur Dibayar',
              value: data.paidInvoicesCount.toString(),
            ),
          ),
        ],
      );

  Widget _buildKpiRow2(BuildContext context) => Row(
        children: [
          Expanded(
            child: SummaryCardWidget(
              icon: Icons.pending_actions_outlined,
              iconColor: AppColors.warningBase,
              iconBackground: AppColors.warningLight,
              label: 'Registrasi Pending',
              value: data.pendingRegistrations.toString(),
            ),
          ),
          const SizedBox(width: AppDimensions.spacing12),
          Expanded(
            child: SummaryCardWidget(
              icon: Icons.person_add_outlined,
              iconColor: AppColors.primary600,
              iconBackground: AppColors.primary50,
              label: 'Registrasi Bulan Ini',
              value: data.newRegistrationsThisMonth.toString(),
            ),
          ),
          const SizedBox(width: AppDimensions.spacing12),
          const Expanded(child: SizedBox()),
        ],
      );

  Widget _buildQuickActions(BuildContext context) => Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Aksi Cepat',
              style: Theme.of(context).textTheme.titleMedium),
          const SizedBox(height: AppDimensions.spacing16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              QuickActionWidget(
                icon: Icons.assignment_outlined,
                iconColor: AppColors.primary600,
                iconBackground: AppColors.primary100,
                label: 'Registrasi',
                onTap: () => context.go('/registrations'),
              ),
              QuickActionWidget(
                icon: Icons.card_membership_outlined,
                iconColor: AppColors.warningBase,
                iconBackground: AppColors.warningLight,
                label: 'Lisensi',
                onTap: () => _showComingSoon(context, 'Lisensi'),
              ),
              QuickActionWidget(
                icon: Icons.policy_outlined,
                iconColor: AppColors.accent400,
                iconBackground: AppColors.accent50,
                label: 'Kebijakan',
                onTap: () => _showComingSoon(context, 'Kebijakan'),
              ),
              QuickActionWidget(
                icon: Icons.settings_outlined,
                iconColor: AppColors.infoBase,
                iconBackground: AppColors.infoLight,
                label: 'Pengaturan',
                onTap: () => _showComingSoon(context, 'Pengaturan Klien'),
              ),
            ],
          ),
        ],
      );

  Widget _buildRecentRegistrations(BuildContext context) => Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Registrasi Terbaru',
                  style: Theme.of(context).textTheme.titleMedium),
              TextButton(
                onPressed: () => context.go('/registrations'),
                child: const Text('Semua'),
              ),
            ],
          ),
          const SizedBox(height: AppDimensions.spacing8),
          ...data.recentRegistrations
              .take(5)
              .map((item) => _buildRegistrationItem(context, item)),
        ],
      );

  Widget _buildRegistrationItem(
      BuildContext context, RecentRegistrationItem item) {
    final (statusColor, statusBg, statusLabel) = switch (item.status) {
      'pending' => (AppColors.warningBase, AppColors.warningLight, 'Pending'),
      'approved' =>
        (AppColors.successBase, AppColors.successLight, 'Disetujui'),
      'rejected' => (AppColors.errorBase, AppColors.errorLight, 'Ditolak'),
      _ => (AppColors.neutral500, AppColors.neutral100, item.status),
    };

    return Container(
      margin: const EdgeInsets.only(bottom: AppDimensions.spacing8),
      padding: const EdgeInsets.all(AppDimensions.cardPadding),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(AppDimensions.radiusLg),
        border: Border.all(color: AppColors.neutral200),
      ),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: AppColors.primary50,
              borderRadius: BorderRadius.circular(AppDimensions.radiusMd),
            ),
            child: const Icon(Icons.business_outlined,
                color: AppColors.primary700, size: AppDimensions.iconMd),
          ),
          const SizedBox(width: AppDimensions.spacing12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(item.companyName,
                    style: Theme.of(context).textTheme.titleSmall),
                const SizedBox(height: AppDimensions.spacing4),
                Text(item.contactName,
                    style: Theme.of(context).textTheme.bodySmall,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(
                horizontal: AppDimensions.spacing8,
                vertical: AppDimensions.spacing4),
            decoration: BoxDecoration(
              color: statusBg,
              borderRadius:
                  BorderRadius.circular(AppDimensions.radiusFull),
            ),
            child: Text(
              statusLabel,
              style: Theme.of(context).textTheme.labelSmall?.copyWith(
                    color: statusColor,
                    fontWeight: FontWeight.w600,
                  ),
            ),
          ),
        ],
      ),
    );
  }

  void _showComingSoon(BuildContext context, String feature) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Fitur $feature belum tersedia')),
    );
  }

  String _labelForRole(String role) {
    switch (role) {
      case 'superuser':
        return 'Superuser';
      case 'finance':
        return 'Finance';
      case 'project_owner':
        return 'Project Owner';
      case 'project_manager':
        return 'Project Manager';
      case 'developer_sales':
        return 'Developer Sales';
      default:
        return role.isNotEmpty ? role : 'Developer';
    }
  }
}
