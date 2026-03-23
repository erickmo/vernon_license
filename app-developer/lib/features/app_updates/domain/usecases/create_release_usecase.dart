import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../repositories/app_update_repository.dart';

class CreateReleaseUseCase {
  final AppUpdateRepository _repo;
  CreateReleaseUseCase(this._repo);

  Future<Either<Failure, void>> call({
    required String appId,
    required String version,
    required int versionCode,
    required String downloadUrl,
    String? releaseNotes,
    bool isMandatory = false,
  }) =>
      _repo.createRelease(
        appId: appId,
        version: version,
        versionCode: versionCode,
        downloadUrl: downloadUrl,
        releaseNotes: releaseNotes,
        isMandatory: isMandatory,
      );
}
