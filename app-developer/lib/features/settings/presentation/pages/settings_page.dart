import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/auth/auth_notifier.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_constants.dart';
import '../../../../injection_container.dart' show sl;

class SettingsPage extends StatelessWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context) {
    final authNotifier = sl<AuthNotifier>();

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Pengaturan'),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // ── Profil ────────────────────────────────────────────────────────
          _SectionHeader(title: 'Profil'),
          _InfoCard(
            children: [
              _InfoTile(
                icon: Icons.badge_outlined,
                label: 'Role',
                value: authNotifier.userRole ?? '-',
              ),
              _InfoTile(
                icon: Icons.fingerprint_rounded,
                label: 'User ID',
                value: authNotifier.userId ?? '-',
                monospace: true,
              ),
            ],
          ),
          const SizedBox(height: 20),

          // ── Server ────────────────────────────────────────────────────────
          _SectionHeader(title: 'Server'),
          _InfoCard(
            children: [
              _InfoTile(
                icon: Icons.dns_outlined,
                label: 'API Base URL',
                value: AppConstants.baseUrl,
                monospace: true,
              ),
              _InfoTile(
                icon: Icons.apps_rounded,
                label: 'Kode Aplikasi',
                value: AppConstants.appCode,
                monospace: true,
              ),
            ],
          ),
          const SizedBox(height: 20),

          // ── Aplikasi ──────────────────────────────────────────────────────
          _SectionHeader(title: 'Aplikasi'),
          _InfoCard(
            children: [
              _InfoTile(
                icon: Icons.info_outline_rounded,
                label: 'Nama Aplikasi',
                value: AppConstants.appName,
              ),
            ],
          ),
          const SizedBox(height: 32),

          // ── Logout ────────────────────────────────────────────────────────
          SizedBox(
            width: double.infinity,
            child: OutlinedButton.icon(
              onPressed: () => _confirmLogout(context, authNotifier),
              icon: const Icon(Icons.logout_rounded),
              label: const Text('Keluar dari Akun'),
              style: OutlinedButton.styleFrom(
                foregroundColor: AppColors.error,
                side: const BorderSide(color: AppColors.error),
                padding: const EdgeInsets.symmetric(vertical: 14),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _confirmLogout(BuildContext context, AuthNotifier authNotifier) {
    showDialog(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text('Keluar dari Akun'),
        content: const Text(
          'Anda akan keluar dari aplikasi. Sesi akan dihapus dan Anda perlu login kembali.',
          style: TextStyle(color: AppColors.neutral700),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Batal'),
          ),
          ElevatedButton(
            onPressed: () async {
              Navigator.pop(context);
              await authNotifier.onLogout();
              if (context.mounted) context.go('/login');
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.error,
            ),
            child: const Text('Keluar'),
          ),
        ],
      ),
    );
  }
}

class _SectionHeader extends StatelessWidget {
  final String title;
  const _SectionHeader({required this.title});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(left: 4, bottom: 8),
      child: Text(
        title.toUpperCase(),
        style: const TextStyle(
          fontSize: 11,
          fontWeight: FontWeight.w700,
          color: AppColors.neutral500,
          letterSpacing: 1.2,
        ),
      ),
    );
  }
}

class _InfoCard extends StatelessWidget {
  final List<Widget> children;
  const _InfoCard({required this.children});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.neutral200),
      ),
      child: Column(
        children: List.generate(children.length, (i) {
          return Column(
            children: [
              children[i],
              if (i < children.length - 1)
                const Divider(height: 1, indent: 48),
            ],
          );
        }),
      ),
    );
  }
}

class _InfoTile extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;
  final bool monospace;

  const _InfoTile({
    required this.icon,
    required this.label,
    required this.value,
    this.monospace = false,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 20, color: AppColors.primary700),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.neutral500,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  value,
                  style: TextStyle(
                    fontSize: 14,
                    color: AppColors.neutral900,
                    fontWeight: FontWeight.w500,
                    fontFamily: monospace ? 'monospace' : null,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
