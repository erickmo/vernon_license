import 'package:flutter_bloc/flutter_bloc.dart';

import '../../domain/entities/registration_entity.dart';
import '../../domain/usecases/approve_registration_usecase.dart';
import '../../domain/usecases/create_client_usecase.dart';
import '../../domain/usecases/list_registrations_usecase.dart';
import '../../domain/usecases/reject_registration_usecase.dart';

part 'registration_state.dart';

class RegistrationCubit extends Cubit<RegistrationState> {
  final ListRegistrationsUseCase _listRegistrations;
  final ApproveRegistrationUseCase _approveRegistration;
  final RejectRegistrationUseCase _rejectRegistration;
  final CreateClientUseCase _createClient;

  String? _activeFilter = 'pending';
  List<RegistrationEntity> _currentItems = [];

  RegistrationCubit({
    required ListRegistrationsUseCase listRegistrations,
    required ApproveRegistrationUseCase approveRegistration,
    required RejectRegistrationUseCase rejectRegistration,
    required CreateClientUseCase createClient,
  })  : _listRegistrations = listRegistrations,
        _approveRegistration = approveRegistration,
        _rejectRegistration = rejectRegistration,
        _createClient = createClient,
        super(RegistrationInitial());

  Future<void> loadRegistrations({String? status}) async {
    _activeFilter = status ?? 'pending';
    emit(RegistrationLoading());

    final result = await _listRegistrations(status: _activeFilter);
    if (result.isLeft()) {
      final failure = result.fold((f) => f, (_) => null)!;
      emit(RegistrationError(failure.message));
      return;
    }

    final items = result.getOrElse(() => []);
    _currentItems = items;
    if (items.isEmpty) {
      emit(RegistrationEmpty(activeFilter: _activeFilter));
    } else {
      emit(RegistrationLoaded(items, activeFilter: _activeFilter));
    }
  }

  Future<void> approve({
    required String id,
    required String companyCode,
    required String companyName,
  }) async {
    emit(RegistrationActionLoading(_currentItems, activeFilter: _activeFilter));

    final result = await _approveRegistration(
      id: id,
      companyCode: companyCode,
      companyName: companyName,
    );

    if (result.isLeft()) {
      final failure = result.fold((f) => f, (_) => null)!;
      emit(RegistrationLoaded(_currentItems, activeFilter: _activeFilter));
      emit(RegistrationError(failure.message));
      return;
    }

    // Reload list setelah approve
    await loadRegistrations(status: _activeFilter);
    emit(RegistrationActionSuccess(
      'Registrasi berhasil disetujui. Perusahaan baru telah dibuat.',
      _currentItems,
      activeFilter: _activeFilter,
    ));
  }

  Future<void> reject({
    required String id,
    required String reason,
  }) async {
    emit(RegistrationActionLoading(_currentItems, activeFilter: _activeFilter));

    final result = await _rejectRegistration(id: id, reason: reason);

    if (result.isLeft()) {
      final failure = result.fold((f) => f, (_) => null)!;
      emit(RegistrationLoaded(_currentItems, activeFilter: _activeFilter));
      emit(RegistrationError(failure.message));
      return;
    }

    await loadRegistrations(status: _activeFilter);
    emit(RegistrationActionSuccess(
      'Registrasi berhasil ditolak.',
      _currentItems,
      activeFilter: _activeFilter,
    ));
  }

  Future<void> createClient({
    required String code,
    required String name,
    required String companyType,
    String? npwp,
    String? email,
    String? phone,
    String? address,
    String? website,
    required List<String> modules,
    required List<String> apps,
  }) async {
    emit(ClientCreateLoading(_currentItems, activeFilter: _activeFilter));

    final result = await _createClient(
      code: code,
      name: name,
      companyType: companyType,
      npwp: npwp,
      email: email,
      phone: phone,
      address: address,
      website: website,
      modules: modules,
      apps: apps,
    );

    if (result.isLeft()) {
      final failure = result.fold((f) => f, (_) => null)!;
      emit(ClientCreateError(failure.message, _currentItems,
          activeFilter: _activeFilter));
      return;
    }

    emit(ClientCreateSuccess(
      'Client "$name" berhasil dibuat.',
      _currentItems,
      activeFilter: _activeFilter,
    ));
  }
}
