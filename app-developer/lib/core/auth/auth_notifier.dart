import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../constants/app_constants.dart';
import '../network/api_client.dart';

class AuthNotifier extends ChangeNotifier {
  final FlutterSecureStorage _storage;
  final ApiClient _apiClient;

  bool _isAuthenticated = false;
  bool get isAuthenticated => _isAuthenticated;

  String? _userRole;
  String? get userRole => _userRole;

  String? _userId;
  String? get userId => _userId;

  AuthNotifier(this._storage, this._apiClient);

  Future<void> init() async {
    final token = await _storage.read(key: AppConstants.accessTokenKey);
    _isAuthenticated = token != null && token.isNotEmpty;
    if (_isAuthenticated) {
      _userRole = await _storage.read(key: AppConstants.userRoleKey);
      _userId = await _storage.read(key: AppConstants.userIdKey);
    }
    notifyListeners();
  }

  Future<void> onLogin(String role, String userId) async {
    await _storage.write(key: AppConstants.userRoleKey, value: role);
    await _storage.write(key: AppConstants.userIdKey, value: userId);
    _userRole = role;
    _userId = userId;
    _isAuthenticated = true;
    notifyListeners();
  }

  Future<void> onLogout() async {
    await _storage.delete(key: AppConstants.accessTokenKey);
    await _storage.delete(key: AppConstants.refreshTokenKey);
    await _storage.delete(key: AppConstants.userRoleKey);
    await _storage.delete(key: AppConstants.userIdKey);
    _isAuthenticated = false;
    _userRole = null;
    _userId = null;
    notifyListeners();
  }

  /// Cek apakah akses app masih valid. Dipanggil secara periodik dan saat app resume.
  /// Jika API mengembalikan 401/403, paksa logout. Error jaringan diabaikan.
  Future<void> checkAccess(String appCode) async {
    if (!_isAuthenticated) return;
    try {
      await _apiClient.dio.get('/api/v1/auth/me');
    } on DioException catch (e) {
      if (e.response?.statusCode == 401 || e.response?.statusCode == 403) {
        await onLogout();
      }
    }
  }
}
