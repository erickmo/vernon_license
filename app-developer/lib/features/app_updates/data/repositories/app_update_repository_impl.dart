import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';

import '../../../../core/errors/failures.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/app_release_entity.dart';
import '../../domain/repositories/app_update_repository.dart';
import '../models/app_release_model.dart';

class AppUpdateRepositoryImpl implements AppUpdateRepository {
  final ApiClient _client;

  AppUpdateRepositoryImpl(this._client);

  @override
  Future<Either<Failure, List<AppReleaseEntity>>> listReleases({String? appId}) async {
    try {
      final queryParams = <String, dynamic>{'limit': 50};
      if (appId != null && appId.isNotEmpty) {
        queryParams['app_id'] = appId;
      }

      final response = await _client.dio.get(
        '/api/v1/developer/app-updates/releases',
        queryParameters: queryParams,
      );

      final items = (response.data['items'] as List? ?? [])
          .map((e) => AppReleaseModel.fromJson(e as Map<String, dynamic>))
          .toList();

      return Right(items);
    } on DioException catch (e) {
      return Left(ServerFailure(e.response?.data?['error'] ?? 'Gagal memuat daftar rilis'));
    }
  }

  @override
  Future<Either<Failure, void>> createRelease({
    required String appId,
    required String version,
    required int versionCode,
    required String downloadUrl,
    String? releaseNotes,
    bool isMandatory = false,
  }) async {
    try {
      await _client.dio.post(
        '/api/v1/developer/app-updates/releases',
        data: {
          'app_id': appId,
          'version': version,
          'version_code': versionCode,
          'download_url': downloadUrl,
          'release_notes': releaseNotes ?? '',
          'is_mandatory': isMandatory,
        },
      );
      return const Right(null);
    } on DioException catch (e) {
      return Left(ServerFailure(e.response?.data?['error'] ?? 'Gagal mempublikasikan rilis'));
    }
  }

  @override
  Future<Either<Failure, void>> pushUpdate({
    required String companyId,
    required String appId,
    required int versionCode,
    bool forceUpdate = false,
  }) async {
    try {
      await _client.dio.post(
        '/api/v1/developer/app-updates/push',
        data: {
          'company_id': companyId,
          'app_id': appId,
          'version_code': versionCode,
          'force_update': forceUpdate,
        },
      );
      return const Right(null);
    } on DioException catch (e) {
      return Left(ServerFailure(e.response?.data?['error'] ?? 'Gagal mendorong update'));
    }
  }

  @override
  Future<Either<Failure, List<ClientInstallEntity>>> getClientInstalls(String companyId) async {
    try {
      final response = await _client.dio.get(
        '/api/v1/developer/app-updates/clients/$companyId/installs',
      );

      final items = (response.data['items'] as List? ?? [])
          .map((e) => ClientInstallModel.fromJson(e as Map<String, dynamic>))
          .toList();

      return Right(items);
    } on DioException catch (e) {
      return Left(ServerFailure(e.response?.data?['error'] ?? 'Gagal memuat data instalasi'));
    }
  }

  @override
  Future<Either<Failure, List<ClientInstallEntity>>> getAppInstalls(String appId) async {
    try {
      final response = await _client.dio.get(
        '/api/v1/developer/app-updates/app/$appId/installs',
      );

      final items = (response.data['items'] as List? ?? [])
          .map((e) => ClientInstallModel.fromJson(e as Map<String, dynamic>))
          .toList();

      return Right(items);
    } on DioException catch (e) {
      return Left(ServerFailure(e.response?.data?['error'] ?? 'Gagal memuat data instalasi'));
    }
  }
}
