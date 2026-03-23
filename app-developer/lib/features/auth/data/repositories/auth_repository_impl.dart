import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/user_entity.dart';
import '../../domain/repositories/auth_repository.dart';

class AuthRepositoryImpl implements AuthRepository {
  final ApiClient _client;
  final FlutterSecureStorage _storage;

  AuthRepositoryImpl(this._client, this._storage);

  @override
  Future<Either<Failure, UserEntity>> login({
    required String identifier,
    required String password,
  }) async {
    try {
      final response = await _client.dio.post(
        '/api/v1/auth/login',
        data: {
          'identifier': identifier,
          'password': password,
        },
      );
      final data = response.data['data'] as Map<String, dynamic>;
      await _storage.write(
          key: AppConstants.accessTokenKey, value: data['access_token'] as String?);

      final user = UserEntity(
        id: data['user_id'] as String? ?? '',
        name: data['name'] as String? ?? '',
        email: data['email'] as String? ?? '',
        role: data['role'] as String? ?? '',
      );

      return Right(user);
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) {
        return const Left(ServerFailure('Email atau password salah'));
      }
      return Left(ServerFailure(
          e.response?.data?['error'] ?? 'Login gagal'));
    } catch (e) {
      return Left(ServerFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> logout() async {
    await _storage.delete(key: AppConstants.accessTokenKey);
    await _storage.delete(key: AppConstants.refreshTokenKey);
    return const Right(null);
  }
}
