class AppConstants {
  AppConstants._();

  static const String appName = 'FlashERP Developer';

  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8081',
  );

  static const Duration connectTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout = Duration(seconds: 30);

  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String userRoleKey = 'user_role';
  static const String userIdKey = 'user_id';

  static const int defaultPageSize = 20;

  /// Kode aplikasi yang dikirim ke API saat login untuk validasi akses.
  static const String appCode = 'app-developer';
}
