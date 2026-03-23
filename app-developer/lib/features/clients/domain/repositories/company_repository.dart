import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../entities/company_entity.dart';

abstract class CompanyRepository {
  Future<Either<Failure, List<CompanyEntity>>> listCompanies({
    String? search,
    bool? activeOnly,
    int limit,
    int offset,
  });
}
