import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';

import '../../../../core/errors/failures.dart';
import '../../domain/entities/developer_dashboard_entity.dart';
import '../../domain/repositories/developer_dashboard_repository.dart';
import '../datasources/developer_dashboard_remote_datasource.dart';

class DeveloperDashboardRepositoryImpl
    implements DeveloperDashboardRepository {
  final DeveloperDashboardRemoteDatasource _datasource;

  DeveloperDashboardRepositoryImpl(this._datasource);

  @override
  Future<Either<Failure, DeveloperDashboardEntity>> getDashboard() async {
    try {
      final result = await _datasource.getDashboard();
      return Right(result);
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) {
        return const Left(UnauthorizedFailure());
      }
      return Left(ServerFailure(
          e.response?.data?['error'] ?? 'Gagal mengambil data dashboard'));
    } catch (_) {
      return const Left(ServerFailure('Gagal mengambil data dashboard'));
    }
  }
}
