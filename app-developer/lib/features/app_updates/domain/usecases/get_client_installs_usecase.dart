import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../entities/app_release_entity.dart';
import '../repositories/app_update_repository.dart';

class GetClientInstallsUseCase {
  final AppUpdateRepository _repo;
  GetClientInstallsUseCase(this._repo);

  Future<Either<Failure, List<ClientInstallEntity>>> call(String companyId) =>
      _repo.getClientInstalls(companyId);
}

class GetAppInstallsUseCase {
  final AppUpdateRepository _repo;
  GetAppInstallsUseCase(this._repo);

  Future<Either<Failure, List<ClientInstallEntity>>> call(String appId) =>
      _repo.getAppInstalls(appId);
}
