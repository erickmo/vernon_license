import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../repositories/registration_repository.dart';

class ApproveRegistrationUseCase {
  final RegistrationRepository _repository;
  const ApproveRegistrationUseCase(this._repository);

  Future<Either<Failure, void>> call({
    required String id,
    required String companyCode,
    required String companyName,
  }) =>
      _repository.approveRegistration(
        id: id,
        companyCode: companyCode,
        companyName: companyName,
      );
}
