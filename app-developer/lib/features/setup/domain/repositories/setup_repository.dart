import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../entities/setup_status_entity.dart';

abstract class SetupRepository {
  Future<Either<Failure, SetupStatusEntity>> getSetupStatus();
  Future<Either<Failure, void>> install({
    required String name,
    required String email,
    required String password,
  });
}
