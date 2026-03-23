import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../entities/developer_dashboard_entity.dart';
import '../repositories/developer_dashboard_repository.dart';

class GetDeveloperDashboardUseCase {
  final DeveloperDashboardRepository _repository;

  GetDeveloperDashboardUseCase(this._repository);

  Future<Either<Failure, DeveloperDashboardEntity>> call() {
    return _repository.getDashboard();
  }
}
