import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';

import '../../../../core/errors/failures.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/registration_entity.dart';
import '../../domain/repositories/registration_repository.dart';
import '../models/registration_model.dart';

class RegistrationRepositoryImpl implements RegistrationRepository {
  final ApiClient _client;

  RegistrationRepositoryImpl(this._client);

  @override
  Future<Either<Failure, List<RegistrationEntity>>> listRegistrations({
    String? status,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page_size': limit,
        'page': (offset ~/ limit) + 1,
      };
      if (status != null && status.isNotEmpty) {
        queryParams['status'] = status;
      }

      final response = await _client.dio.get(
        '/api/v1/licenses',
        queryParameters: queryParams,
      );

      final body = response.data['data'] as Map<String, dynamic>? ?? {};
      final items = (body['items'] as List? ?? [])
          .map((e) => RegistrationModel.fromJson(e as Map<String, dynamic>))
          .toList();

      return Right(items);
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) {
        return const Left(UnauthorizedFailure());
      }
      return Left(ServerFailure(
          e.response?.data?['error'] ?? 'Gagal mengambil daftar lisensi'));
    } catch (e) {
      return Left(ServerFailure('Gagal mengambil daftar lisensi'));
    }
  }

  @override
  Future<Either<Failure, void>> approveRegistration({
    required String id,
    required String companyCode,
    required String companyName,
  }) async {
    try {
      await _client.dio.put(
        '/api/v1/licenses/$id/status',
        data: {'status': 'active'},
      );
      return const Right(null);
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) {
        return const Left(UnauthorizedFailure());
      }
      return Left(ServerFailure(
          e.response?.data?['error'] ?? 'Gagal mengaktifkan lisensi'));
    } catch (e) {
      return Left(ServerFailure('Gagal mengaktifkan lisensi'));
    }
  }

  @override
  Future<Either<Failure, void>> rejectRegistration({
    required String id,
    required String reason,
  }) async {
    try {
      await _client.dio.put(
        '/api/v1/licenses/$id/status',
        data: {'status': 'suspended'},
      );
      return const Right(null);
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) {
        return const Left(UnauthorizedFailure());
      }
      return Left(ServerFailure(
          e.response?.data?['error'] ?? 'Gagal menangguhkan lisensi'));
    } catch (e) {
      return Left(ServerFailure('Gagal menangguhkan lisensi'));
    }
  }

  @override
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
  }) async {
    try {
      final body = <String, dynamic>{
        'client_name': name,
        'client_email': email ?? '',
        'plan': companyType,
        'product': 'flasherp',
      };
      if (website != null && website.isNotEmpty) body['flasherp_url'] = website;

      await _client.dio.post('/api/v1/licenses', data: body);
      return const Right(null);
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) {
        return const Left(UnauthorizedFailure());
      }
      return Left(ServerFailure(
          e.response?.data?['error'] ?? 'Gagal membuat lisensi baru'));
    } catch (e) {
      return Left(ServerFailure('Gagal membuat lisensi baru'));
    }
  }
}
