import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../repositories/registration_repository.dart';

class CreateClientUseCase {
  final RegistrationRepository _repository;

  CreateClientUseCase(this._repository);

  Future<Either<Failure, void>> call({
    required String code,
    required String name,
    required String companyType,
    String? npwp,
    String? email,
    String? phone,
    String? address,
    String? website,
    required List<String> modules,
    required List<String> apps,
  }) =>
      _repository.createClient(
        code: code,
        name: name,
        companyType: companyType,
        npwp: npwp,
        email: email,
        phone: phone,
        address: address,
        website: website,
        modules: modules,
        apps: apps,
      );
}
