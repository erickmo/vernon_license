import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'core/auth/auth_notifier.dart';
import 'core/constants/app_constants.dart';
import 'core/setup/setup_notifier.dart';
import 'core/theme/app_theme.dart';
import 'core/widgets/main_shell.dart';
import 'features/auth/presentation/pages/login_page.dart';
import 'features/dashboard/presentation/pages/dashboard_page.dart';
import 'features/invoices/presentation/pages/invoice_page.dart';
import 'features/notifications/presentation/pages/notifications_page.dart';
import 'features/clients/presentation/pages/clients_page.dart';
import 'features/settings/presentation/pages/settings_page.dart';
import 'features/setup/presentation/pages/setup_page.dart';
import 'features/app_updates/presentation/pages/app_updates_page.dart';
import 'injection_container.dart' show sl;

class App extends StatefulWidget {
  const App({super.key});

  @override
  State<App> createState() => _AppState();
}

class _AppState extends State<App> with WidgetsBindingObserver {
  Timer? _accessCheckTimer;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    WidgetsBinding.instance.addPostFrameCallback((_) => _startTimer());
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    _accessCheckTimer?.cancel();
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (state == AppLifecycleState.resumed) {
      sl<AuthNotifier>().checkAccess(AppConstants.appCode);
    }
  }

  void _startTimer() {
    _accessCheckTimer = Timer.periodic(const Duration(minutes: 5), (_) {
      sl<AuthNotifier>().checkAccess(AppConstants.appCode);
    });
  }

  @override
  Widget build(BuildContext context) {
    final authNotifier = sl<AuthNotifier>();
    final setupNotifier = sl<SetupNotifier>();

    final router = GoRouter(
      initialLocation: '/login',
      refreshListenable: Listenable.merge([authNotifier, setupNotifier]),
      redirect: (context, state) {
        final loc = state.matchedLocation;

        if (setupNotifier.isChecking) return null;

        if (!setupNotifier.isInstalled && loc != '/setup') return '/setup';
        if (setupNotifier.isInstalled && loc == '/setup') return '/login';

        if (!authNotifier.isAuthenticated && loc != '/login') return '/login';
        if (authNotifier.isAuthenticated && loc == '/login') return '/home';

        return null;
      },
      routes: [
        GoRoute(
          path: '/login',
          builder: (_, __) => const LoginPage(),
        ),
        GoRoute(
          path: '/setup',
          builder: (_, __) => const SetupPage(),
        ),
        ShellRoute(
          builder: (context, state, child) => MainShell(
            location: state.uri.path,
            child: child,
          ),
          routes: [
            GoRoute(
              path: '/home',
              builder: (_, __) => const DashboardPage(),
            ),
            GoRoute(
              path: '/clients',
              builder: (_, __) => const ClientsPage(),
            ),
            GoRoute(
              path: '/notifications',
              builder: (_, __) => const NotificationsPage(),
            ),
            GoRoute(
              path: '/invoices',
              builder: (_, __) => const InvoicePage(),
            ),
            GoRoute(
              path: '/updates',
              builder: (_, __) => const AppUpdatesPage(),
            ),
            GoRoute(
              path: '/settings',
              builder: (_, __) => const SettingsPage(),
            ),
          ],
        ),
      ],
    );

    return MaterialApp.router(
      title: 'FlashERP Developer',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light,
      routerConfig: router,
    );
  }
}
