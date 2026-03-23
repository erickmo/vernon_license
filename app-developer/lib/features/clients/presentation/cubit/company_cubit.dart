import 'package:flutter_bloc/flutter_bloc.dart';

import '../../domain/entities/company_entity.dart';
import '../../domain/usecases/list_companies_usecase.dart';

part 'company_state.dart';

class CompanyCubit extends Cubit<CompanyState> {
  final ListCompaniesUseCase _listCompanies;

  String _search = '';

  CompanyCubit(this._listCompanies) : super(CompanyInitial());

  Future<void> load({String? search}) async {
    _search = search ?? '';
    emit(CompanyLoading());

    final result = await _listCompanies(search: _search.isEmpty ? null : _search);

    result.fold(
      (failure) => emit(CompanyError(failure.message)),
      (items) {
        if (items.isEmpty) {
          emit(CompanyEmpty(search: _search));
        } else {
          emit(CompanyLoaded(items, search: _search));
        }
      },
    );
  }

  void search(String query) => load(search: query);

  Future<void> refresh() => load(search: _search);
}
