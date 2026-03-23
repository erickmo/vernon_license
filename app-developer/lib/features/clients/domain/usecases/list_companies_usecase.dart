import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../entities/company_entity.dart';
import '../repositories/company_repository.dart';

class ListCompaniesUseCase {
  final CompanyRepository _repository;

  ListCompaniesUseCase(this._repository);

  Future<Either<Failure, List<CompanyEntity>>> call({
    String? search,
    bool? activeOnly,
    int limit = 50,
    int offset = 0,
  }) =>
      _repository.listCompanies(
        search: search,
        activeOnly: activeOnly,
        limit: limit,
        offset: offset,
      );
}
