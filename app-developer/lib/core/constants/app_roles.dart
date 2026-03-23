class AppRoles {
  AppRoles._();

  static const developerSales = 'developer_sales';

  static const allowedRoles = [developerSales];

  static bool isAllowed(String role) => allowedRoles.contains(role);
}
