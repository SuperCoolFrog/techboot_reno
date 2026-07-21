% :- debug.
% Data

% .State Scene 3
state(3).

file_details(file_id(1), file_type(email), file_name(roy_1)).
file_details(file_id(2), file_type(email), file_name(roy_2)).
file_details(file_id(3), file_type(email), file_name(roy_3)).


fs(state(3),file_id(1)). 
fs(state(3),file_id(2)). 
fs(state(3),file_id(3)). 

shape(1, square).

connection(state(3), rabbit).

result(invalid).
result(connect_true).
result(connect_false).

print_results(A, B) :- format('The result is: ~w ;; ~w~n', [A, B]).

print_result_out(_, B) :- 
    out_value(B, Out),
    format('The Out is: ~w~n', [Out]).


out_value([_|Rest], Out) :- Rest \= [], !, out_value(Rest, Out).
out_value([H|_], Out) :- Out = H.


% Predicates

search(FT, FN, StateId) :- 
    state(StateId), 
    fs(state(StateId), file_id(FID)), 
    file_details(file_id(FID), file_type(FT), file_name(FN)).

connect(Name, StateId, Out) :-
		state(StateId),
		connection(state(StateId), Name), !,
        Out = result(connect_true).

connect(_, _, Out) :- Out = result(connect_false).

% Parser
process_command(List, Out) :-
    % 1. Clean filler words
    remove_fillers(List, CleanedList),
    
    % 2. Extract Verb and raw text arguments
    CleanedList = [Verb | RawArgs],
    
    % 3. Convert text placeholders (like 'FT') into true Prolog variables
    maplist(bind_variables, RawArgs, PrologArgs),
    
    % 4. Construct the executable Goal
    Goal =.. [Verb | PrologArgs],
    
    % 5. ISO-compliant safety check: extract Name/Arity and check existence
    functor(Goal, Name, Arity),
    current_predicate(Name/Arity),
    call(Goal),
    
    % 6. Print the answers to the user console
    % print_results(RawArgs, PrologArgs),
    % print_result_out(RawArgs, PrologArgs),
    !, out_value(PrologArgs, Out).

% Any commands that fail return invalid
process_command(_, Out) :- Out = result(invalid).

% Helper: If an atom starts with an uppercase letter or is an underscore, treat it as a variable
bind_variables(Atom, _) :-
    atom(Atom),
    (   sub_atom(Atom, 0, 1, _, FirstChar),
        char_code(FirstChar, Code),
        Code >= 65, Code =< 90 % Matches 'A' through 'Z'
    ;   Atom == '_'
    ),
    !. % If it matches the criteria, leave Variable unbound.
bind_variables(Value, Value). % Otherwise, keep it as a concrete value (like 'my_file.txt').

remove_fillers([], []).
remove_fillers([X|Xs], Ys) :-
    member(X, [to, from, at, with, the, a]), !, % Add any filler words here
    remove_fillers(Xs, Ys).
remove_fillers([X|Xs], [X|Ys]) :-
    remove_fillers(Xs, Ys).

% Examples
% process_command([connect, rabbit, 3, 'Out'], Out).

