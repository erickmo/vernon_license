import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../constants/app_colors.dart';
import '../notifiers/pending_count_notifier.dart';
import '../../injection_container.dart' show sl;

/// Shell utama aplikasi — membungkus semua halaman authenticated
/// dengan [NavigationBar] 5 tab: Beranda, Client, Invoice, Notifikasi, Pengaturan.
class MainShell extends StatelessWidget {
  final Widget child;
  final String location;

  const MainShell({
    super.key,
    required this.child,
    required this.location,
  });

  static const _routes = [
    '/home',
    '/clients',
    '/updates',
    '/notifications',
    '/settings',
  ];

  static int _locationToIndex(String loc) {
    if (loc.startsWith('/clients')) return 1;
    if (loc.startsWith('/updates')) return 2;
    if (loc.startsWith('/notifications')) return 3;
    if (loc.startsWith('/settings')) return 4;
    return 0;
  }

  @override
  Widget build(BuildContext context) {
    final selectedIndex = _locationToIndex(location);
    final pendingNotifier = sl<PendingCountNotifier>();

    return Scaffold(
      body: child,
      bottomNavigationBar: ValueListenableBuilder<int>(
        valueListenable: pendingNotifier,
        builder: (_, count, __) => NavigationBar(
          selectedIndex: selectedIndex,
          backgroundColor: AppColors.surface,
          indicatorColor: AppColors.primary100,
          labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
          onDestinationSelected: (index) => context.go(_routes[index]),
          destinations: [
            const NavigationDestination(
              icon: Icon(Icons.home_outlined),
              selectedIcon: Icon(Icons.home_rounded),
              label: 'Beranda',
            ),
            const NavigationDestination(
              icon: Icon(Icons.people_outline_rounded),
              selectedIcon: Icon(Icons.people_rounded),
              label: 'Client',
            ),
            const NavigationDestination(
              icon: Icon(Icons.system_update_outlined),
              selectedIcon: Icon(Icons.system_update_rounded),
              label: 'Update App',
            ),
            NavigationDestination(
              icon: Badge(
                isLabelVisible: count > 0,
                label: Text('$count'),
                child: const Icon(Icons.notifications_outlined),
              ),
              selectedIcon: Badge(
                isLabelVisible: count > 0,
                label: Text('$count'),
                child: const Icon(Icons.notifications_rounded),
              ),
              label: 'Notifikasi',
            ),
            const NavigationDestination(
              icon: Icon(Icons.settings_outlined),
              selectedIcon: Icon(Icons.settings_rounded),
              label: 'Pengaturan',
            ),
          ],
        ),
      ),
    );
  }
}
