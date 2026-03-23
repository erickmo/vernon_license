part of 'developer_dashboard_cubit.dart';

abstract class DeveloperDashboardState extends Equatable {
  const DeveloperDashboardState();
  @override
  List<Object?> get props => [];
}

class DeveloperDashboardInitial extends DeveloperDashboardState {}

class DeveloperDashboardLoading extends DeveloperDashboardState {}

class DeveloperDashboardLoaded extends DeveloperDashboardState {
  final DeveloperDashboardEntity data;
  const DeveloperDashboardLoaded(this.data);
  @override
  List<Object?> get props => [data];
}

class DeveloperDashboardError extends DeveloperDashboardState {
  final String message;
  const DeveloperDashboardError(this.message);
  @override
  List<Object?> get props => [message];
}
