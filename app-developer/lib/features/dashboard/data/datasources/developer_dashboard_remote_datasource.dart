import '../../../../core/network/api_client.dart';
import '../models/developer_dashboard_model.dart';

abstract class DeveloperDashboardRemoteDatasource {
  Future<DeveloperDashboardModel> getDashboard();
}

class DeveloperDashboardRemoteDatasourceImpl
    implements DeveloperDashboardRemoteDatasource {
  final ApiClient _client;

  DeveloperDashboardRemoteDatasourceImpl(this._client);

  @override
  Future<DeveloperDashboardModel> getDashboard() async {
    final response = await _client.dio.get('/api/v1/dashboard');
    final data = response.data['data'] as Map<String, dynamic>? ?? {};
    return DeveloperDashboardModel.fromJson(data);
  }
}
