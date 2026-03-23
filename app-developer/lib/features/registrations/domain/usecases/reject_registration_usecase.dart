import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../repositories/registration_repository.dart';

class RejectRegistrationUseCase {
  final RegistrationRepository _repository;
  const RejectRegistrationUseCase(this._repository);

  Future<Either<Failure, void>> call({
    required String id,
    required String reason,
  }) =>
      _repository.rejectRegistration(id: id, reason: reason);
}
