import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';

import '../../../../core/errors/failures.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/setup_status_entity.dart';
import '../../domain/repositories/setup_repository.dart';

class SetupRepositoryImpl implements SetupRepository {
  final ApiClient _client;

  SetupRepositoryImpl(this._client);

  @override
  Future<Either<Failure, SetupStatusEntity>> getSetupStatus() async {
    try {
      final response = await _client.dio.get('/api/v1/setup/status');
      final data = response.data['data'];
      return Right(SetupStatusEntity(
        isInstalled: data['is_installed'] as bool? ?? false,
      ));
    } on DioException catch (e) {
      return Left(ServerFailure(
          e.response?.data?['error'] ?? 'Gagal mengambil status setup'));
    } catch (e) {
      return const Left(ServerFailure('Gagal mengambil status setup'));
    }
  }

  @override
  Future<Either<Failure, void>> install({
    required String name,
    required String email,
    required String password,
  }) async {
    try {
      await _client.dio.post(
        '/api/v1/setup/install',
        data: {
          'name': name,
          'email': email,
          'password': password,
        },
      );
      return const Right(null);
    } on DioException catch (e) {
      return Left(ServerFailure(
          e.response?.data?['error'] ?? 'Gagal melakukan instalasi'));
    } catch (e) {
      return const Left(ServerFailure('Gagal melakukan instalasi'));
    }
  }
}
