import '../../domain/entities/developer_dashboard_entity.dart';

class DeveloperDashboardModel extends DeveloperDashboardEntity {
  const DeveloperDashboardModel({
    required super.totalClients,
    required super.activeCompanies,
    required super.paidInvoicesCount,
    required super.mrr,
    required super.mrrGrowthPercent,
    required super.pendingRegistrations,
    required super.newRegistrationsThisMonth,
    required super.recentRegistrations,
  });

  factory DeveloperDashboardModel.fromJson(Map<String, dynamic> json) {
    return DeveloperDashboardModel(
      totalClients: (json['total_clients'] as num? ?? 0).toInt(),
      activeCompanies: (json['active_companies'] as num? ?? 0).toInt(),
      paidInvoicesCount: (json['paid_invoices_count'] as num? ?? 0).toInt(),
      mrr: (json['mrr'] as num? ?? 0).toDouble(),
      mrrGrowthPercent: (json['mrr_growth_percent'] as num? ?? 0).toDouble(),
      pendingRegistrations:
          (json['pending_registrations'] as num? ?? 0).toInt(),
      newRegistrationsThisMonth:
          (json['new_registrations_this_month'] as num? ?? 0).toInt(),
      recentRegistrations:
          ((json['recent_registrations'] as List?) ?? [])
              .map((e) =>
                  RecentRegistrationModel.fromJson(e as Map<String, dynamic>))
              .toList(),
    );
  }
}

class RecentRegistrationModel extends RecentRegistrationItem {
  const RecentRegistrationModel({
    required super.id,
    required super.companyName,
    required super.contactName,
    required super.status,
    required super.createdAt,
  });

  factory RecentRegistrationModel.fromJson(Map<String, dynamic> json) {
    return RecentRegistrationModel(
      id: json['id'] as String? ?? '',
      companyName: json['company_name'] as String? ?? '',
      contactName: json['contact_name'] as String? ?? '',
      status: json['status'] as String? ?? '',
      createdAt: DateTime.tryParse(json['created_at'] as String? ?? '') ??
          DateTime.now(),
    );
  }
}
