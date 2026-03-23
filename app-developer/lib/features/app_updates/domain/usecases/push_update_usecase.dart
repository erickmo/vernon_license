import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../repositories/app_update_repository.dart';

class PushUpdateUseCase {
  final AppUpdateRepository _repo;
  PushUpdateUseCase(this._repo);

  Future<Either<Failure, void>> call({
    required String companyId,
    required String appId,
    required int versionCode,
    bool forceUpdate = false,
  }) =>
      _repo.pushUpdate(
        companyId: companyId,
        appId: appId,
        versionCode: versionCode,
        forceUpdate: forceUpdate,
      );
}
