import 'package:get_it/get_it.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import 'core/auth/auth_notifier.dart';
import 'core/network/api_client.dart';
import 'core/notifiers/pending_count_notifier.dart';
import 'core/setup/setup_notifier.dart';

import 'features/auth/data/repositories/auth_repository_impl.dart';
import 'features/auth/domain/repositories/auth_repository.dart';
import 'features/auth/domain/usecases/login_usecase.dart';
import 'features/auth/presentation/cubit/auth_cubit.dart';

import 'features/dashboard/data/datasources/developer_dashboard_remote_datasource.dart';
import 'features/dashboard/data/repositories/developer_dashboard_repository_impl.dart';
import 'features/dashboard/domain/repositories/developer_dashboard_repository.dart';
import 'features/dashboard/domain/usecases/get_developer_dashboard_usecase.dart';
import 'features/dashboard/presentation/cubit/developer_dashboard_cubit.dart';

import 'features/clients/data/repositories/company_repository_impl.dart';
import 'features/clients/domain/repositories/company_repository.dart';
import 'features/clients/domain/usecases/list_companies_usecase.dart';
import 'features/clients/presentation/cubit/company_cubit.dart';

import 'features/registrations/data/repositories/registration_repository_impl.dart';
import 'features/registrations/domain/repositories/registration_repository.dart';
import 'features/registrations/domain/usecases/approve_registration_usecase.dart';
import 'features/registrations/domain/usecases/create_client_usecase.dart';
import 'features/registrations/domain/usecases/list_registrations_usecase.dart';
import 'features/registrations/domain/usecases/reject_registration_usecase.dart';
import 'features/registrations/presentation/cubit/registration_cubit.dart';

import 'features/setup/data/repositories/setup_repository_impl.dart';
import 'features/setup/domain/repositories/setup_repository.dart';
import 'features/setup/presentation/cubit/setup_cubit.dart';

import 'features/app_updates/data/repositories/app_update_repository_impl.dart';
import 'features/app_updates/domain/repositories/app_update_repository.dart';
import 'features/app_updates/domain/usecases/create_release_usecase.dart';
import 'features/app_updates/domain/usecases/get_client_installs_usecase.dart';
import 'features/app_updates/domain/usecases/list_releases_usecase.dart';
import 'features/app_updates/domain/usecases/push_update_usecase.dart';
import 'features/app_updates/presentation/cubit/app_update_cubit.dart';

final sl = GetIt.instance;

Future<void> init() async {
  // Core
  sl.registerLazySingleton(() => const FlutterSecureStorage());
  sl.registerLazySingleton(() => ApiClient(sl()));
  sl.registerLazySingleton(() => AuthNotifier(sl(), sl()));
  sl.registerLazySingleton(() => SetupNotifier());
  sl.registerLazySingleton(() => PendingCountNotifier());

  // Auth
  sl.registerLazySingleton<AuthRepository>(() => AuthRepositoryImpl(sl(), sl()));
  sl.registerLazySingleton(() => LoginUseCase(sl()));
  sl.registerFactory(() => AuthCubit(loginUseCase: sl(), authNotifier: sl()));

  // Dashboard
  sl.registerLazySingleton<DeveloperDashboardRemoteDatasource>(
      () => DeveloperDashboardRemoteDatasourceImpl(sl()));
  sl.registerLazySingleton<DeveloperDashboardRepository>(
      () => DeveloperDashboardRepositoryImpl(sl()));
  sl.registerLazySingleton(() => GetDeveloperDashboardUseCase(sl()));
  sl.registerFactory(() => DeveloperDashboardCubit(sl()));

  // Clients (Companies)
  sl.registerLazySingleton<CompanyRepository>(() => CompanyRepositoryImpl(sl()));
  sl.registerLazySingleton(() => ListCompaniesUseCase(sl()));
  sl.registerFactory(() => CompanyCubit(sl()));

  // Registrations
  sl.registerLazySingleton<RegistrationRepository>(
      () => RegistrationRepositoryImpl(sl()));
  sl.registerLazySingleton(() => ListRegistrationsUseCase(sl()));
  sl.registerLazySingleton(() => ApproveRegistrationUseCase(sl()));
  sl.registerLazySingleton(() => RejectRegistrationUseCase(sl()));
  sl.registerLazySingleton(() => CreateClientUseCase(sl()));
  sl.registerFactory(() => RegistrationCubit(
        listRegistrations: sl(),
        approveRegistration: sl(),
        rejectRegistration: sl(),
        createClient: sl(),
      ));

  // Setup
  sl.registerLazySingleton<SetupRepository>(() => SetupRepositoryImpl(sl()));
  sl.registerFactory(() => SetupCubit(setupRepository: sl()));

  // App Updates (OTA)
  sl.registerLazySingleton<AppUpdateRepository>(() => AppUpdateRepositoryImpl(sl()));
  sl.registerLazySingleton(() => ListReleasesUseCase(sl()));
  sl.registerLazySingleton(() => CreateReleaseUseCase(sl()));
  sl.registerLazySingleton(() => PushUpdateUseCase(sl()));
  sl.registerLazySingleton(() => GetClientInstallsUseCase(sl()));
  sl.registerLazySingleton(() => GetAppInstallsUseCase(sl()));
  sl.registerFactory(() => AppUpdateCubit(
        listReleases: sl(),
        createRelease: sl(),
        pushUpdate: sl(),
      ));
}
