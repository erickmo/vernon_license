import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../entities/registration_entity.dart';

abstract class RegistrationRepository {
  Future<Either<Failure, List<RegistrationEntity>>> listRegistrations({
    String? status,
    int limit,
    int offset,
  });

  Future<Either<Failure, void>> approveRegistration({
    required String id,
    required String companyCode,
    required String companyName,
  });

  Future<Either<Failure, void>> rejectRegistration({
    required String id,
    required String reason,
  });

  Future<Either<Failure, void>> createClient({
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
  });
}
