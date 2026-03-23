import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../entities/registration_entity.dart';
import '../repositories/registration_repository.dart';

class ListRegistrationsUseCase {
  final RegistrationRepository _repository;
  const ListRegistrationsUseCase(this._repository);

  Future<Either<Failure, List<RegistrationEntity>>> call({
    String? status,
    int limit = 20,
    int offset = 0,
  }) =>
      _repository.listRegistrations(
          status: status, limit: limit, offset: offset);
}
