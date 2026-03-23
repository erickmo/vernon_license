import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';

import '../../../../core/errors/failures.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/company_entity.dart';
import '../../domain/repositories/company_repository.dart';
import '../models/company_model.dart';

class CompanyRepositoryImpl implements CompanyRepository {
  final ApiClient _client;

  CompanyRepositoryImpl(this._client);

  @override
  Future<Either<Failure, List<CompanyEntity>>> listCompanies({
    String? search,
    bool? activeOnly,
    int limit = 50,
    int offset = 0,
  }) async {
    try {
      final params = <String, dynamic>{'limit': limit, 'offset': offset};
      if (search != null && search.isNotEmpty) params['search'] = search;
      if (activeOnly != null) params['active_only'] = activeOnly;

      final response = await _client.dio.get(
        '/api/v1/licenses',
        queryParameters: params,
      );

      final body = response.data['data'] as Map<String, dynamic>? ?? {};
      final rawList = (body['items'] ?? []) as List;
      final items = rawList
          .map((e) => CompanyModel.fromJson(e as Map<String, dynamic>))
          .toList();

      return Right(items);
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) return const Left(UnauthorizedFailure());
      return Left(ServerFailure(
          e.response?.data?['error'] ?? 'Gagal mengambil daftar client'));
    } catch (_) {
      return const Left(ServerFailure('Gagal mengambil daftar client'));
    }
  }
}
