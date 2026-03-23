import 'package:flutter_bloc/flutter_bloc.dart';

import '../../domain/repositories/setup_repository.dart';

part 'setup_state.dart';

class SetupCubit extends Cubit<SetupState> {
  final SetupRepository _setupRepository;

  SetupCubit({required SetupRepository setupRepository})
      : _setupRepository = setupRepository,
        super(SetupInitial());

  Future<void> checkStatus() async {
    emit(SetupLoading());
    final result = await _setupRepository.getSetupStatus();
    result.fold(
      (failure) => emit(SetupError(failure.message)),
      (status) => emit(SetupStatusLoaded(status.isInstalled)),
    );
  }

  Future<void> install({
    required String name,
    required String email,
    required String password,
  }) async {
    emit(SetupInstalling());
    final result = await _setupRepository.install(
      name: name,
      email: email,
      password: password,
    );
    result.fold(
      (failure) => emit(SetupError(failure.message)),
      (_) => emit(SetupInstallSuccess()),
    );
  }
}
