import 'package:equatable/equatable.dart';

class DeveloperDashboardEntity extends Equatable {
  final int totalClients;
  final int activeCompanies;
  final int paidInvoicesCount;
  final double mrr;
  final double mrrGrowthPercent;
  final int pendingRegistrations;
  final int newRegistrationsThisMonth;
  final List<RecentRegistrationItem> recentRegistrations;

  const DeveloperDashboardEntity({
    required this.totalClients,
    required this.activeCompanies,
    required this.paidInvoicesCount,
    required this.mrr,
    required this.mrrGrowthPercent,
    required this.pendingRegistrations,
    required this.newRegistrationsThisMonth,
    required this.recentRegistrations,
  });

  @override
  List<Object> get props => [
        totalClients,
        activeCompanies,
        paidInvoicesCount,
        mrr,
        mrrGrowthPercent,
        pendingRegistrations,
        newRegistrationsThisMonth,
        recentRegistrations,
      ];
}

class RecentRegistrationItem extends Equatable {
  final String id;
  final String companyName;
  final String contactName;
  final String status;
  final DateTime createdAt;

  const RecentRegistrationItem({
    required this.id,
    required this.companyName,
    required this.contactName,
    required this.status,
    required this.createdAt,
  });

  @override
  List<Object> get props => [id, companyName, contactName, status, createdAt];
}
