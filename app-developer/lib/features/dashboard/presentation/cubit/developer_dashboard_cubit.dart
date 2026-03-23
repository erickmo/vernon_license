import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../domain/entities/developer_dashboard_entity.dart';
import '../../domain/usecases/get_developer_dashboard_usecase.dart';

part 'developer_dashboard_state.dart';

class DeveloperDashboardCubit extends Cubit<DeveloperDashboardState> {
  final GetDeveloperDashboardUseCase _useCase;

  DeveloperDashboardCubit(this._useCase) : super(DeveloperDashboardInitial());

  Future<void> loadDashboard() async {
    emit(DeveloperDashboardLoading());
    final result = await _useCase();
    result.fold(
      (failure) => emit(DeveloperDashboardError(failure.message)),
      (data) => emit(DeveloperDashboardLoaded(data)),
    );
  }
}
