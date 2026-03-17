import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '../../../models/anime.dart';
import '../../../core/bridge.dart';

// Events
abstract class SearchEvent extends Equatable {
  @override
  List<Object> get props => [];
}

class SearchQueryChanged extends SearchEvent {
  final String query;
  SearchQueryChanged(this.query);
  @override
  List<Object> get props => [query];
}

// States
abstract class SearchState extends Equatable {
  @override
  List<Object> get props => [];
}

class SearchInitial extends SearchState {}
class SearchLoading extends SearchState {}
class SearchSuccess extends SearchState {
  final List<Anime> results;
  SearchSuccess(this.results);
  @override
  List<Object> get props => [results];
}
class SearchError extends SearchState {
  final String message;
  SearchError(this.message);
  @override
  List<Object> get props => [message];
}

// BLoC
class SearchBloc extends Bloc<SearchEvent, SearchState> {
  SearchBloc() : super(SearchInitial()) {
    on<SearchQueryChanged>((event, emit) async {
      if (event.query.isEmpty) {
        emit(SearchInitial());
        return;
      }
      
      emit(SearchLoading());
      try {
        final List<dynamic> data = await GoBridge.search(event.query);
        final List<Anime> results = data.map((e) => Anime.fromJson(e)).toList();
        emit(SearchSuccess(results));
      } catch (e) {
        emit(SearchError(e.toString()));
      }
    });
  }
}
