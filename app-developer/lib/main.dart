import 'package:flutter/material.dart';
import 'package:intl/date_symbol_data_local.dart';

import 'app.dart';
import 'injection_container.dart' as di;
import 'core/auth/auth_notifier.dart';
import 'core/setup/setup_notifier.dart';
import 'features/setup/domain/repositories/setup_repository.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await initializeDateFormatting('id_ID', null);
  await di.init();
  await di.sl<AuthNotifier>().init();

  // Cek setup status
  final setupNotifier = di.sl<SetupNotifier>();
  final setupRepo = di.sl<SetupRepository>();
  setupRepo.getSetupStatus().then((result) {
    result.fold(
      (failure) => setupNotifier.setInstalled(true), // fallback: anggap installed jika error
      (status) => setupNotifier.setInstalled(status.isInstalled),
    );
  });

  runApp(const App());
}
