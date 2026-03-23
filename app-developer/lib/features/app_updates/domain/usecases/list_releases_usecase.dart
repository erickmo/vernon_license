import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../entities/app_release_entity.dart';
import '../repositories/app_update_repository.dart';

class ListReleasesUseCase {
  final AppUpdateRepository _repo;
  ListReleasesUseCase(this._repo);

  Future<Either<Failure, List<AppReleaseEntity>>> call({String? appId}) =>
      _repo.listReleases(appId: appId);
}
