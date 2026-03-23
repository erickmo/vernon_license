import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../entities/app_release_entity.dart';

abstract class AppUpdateRepository {
  Future<Either<Failure, List<AppReleaseEntity>>> listReleases({String? appId});

  Future<Either<Failure, void>> createRelease({
    required String appId,
    required String version,
    required int versionCode,
    required String downloadUrl,
    String? releaseNotes,
    bool isMandatory = false,
  });

  Future<Either<Failure, void>> pushUpdate({
    required String companyId,
    required String appId,
    required int versionCode,
    bool forceUpdate = false,
  });

  Future<Either<Failure, List<ClientInstallEntity>>> getClientInstalls(String companyId);

  Future<Either<Failure, List<ClientInstallEntity>>> getAppInstalls(String appId);
}
