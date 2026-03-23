import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../entities/developer_dashboard_entity.dart';

abstract class DeveloperDashboardRepository {
  Future<Either<Failure, DeveloperDashboardEntity>> getDashboard();
}
