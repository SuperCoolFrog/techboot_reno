% :- debug.
% Data

% .State Scene 3
state(21).

file_details(file_id(1), file_type(email), file_name(roy_1)).
file_details(file_id(2), file_type(email), file_name(roy_2)).
file_details(file_id(3), file_type(email), file_name(roy_3)).

fs(state(999),file_id(1)). 
fs(state(999),file_id(2)). 
fs(state(999),file_id(3)). 

shape(1, square).


% Scene3_InputHandlingLoop
connection(state(21), rabbit).

result(invalid).
result(connect_true).
result(connect_false).

print_results(A, B) :- format('The result is: ~w ;; ~w~n', [A, B]).

print_result_out(_, B) :- 
    out_value(B, Out),
    format('The Out is: ~w~n', [Out]).

out_value([_|Rest], Out) :- Rest \= [], !, out_value(Rest, Out).
out_value([H|_], Out) :- Out = H.

% =========================================================================
% Target Predicates (Refactored for Clean Unification)
% =========================================================================
search(FT, FN, StateId) :- 
    state(StateId), 
    fs(state(StateId), file_id(FID)), 
    file_details(file_id(FID), file_type(FT), file_name(FN)).

% Pattern-match result directly in the head to force early unification
connect(Name, StateId, result(connect_true)) :-
    state(StateId),
    connection(state(StateId), Name), !.

connect(_, _, result(connect_false)).

% =========================================================================
% Parser
% =========================================================================
process_command(List, Out) :-
    % 1. Clean filler words
    remove_fillers(List, CleanedList),
    
    % 2. Extract Verb and raw text arguments
    CleanedList = [Verb | RawArgs],
    
    % 3. Convert text placeholders (like 'Out') into true Prolog variables
    maplist(bind_variables, RawArgs, PrologArgs),
    
    % 4. Unify the outer 'Out' variable with the list parameter early 
    % This locks the memory addresses together for Trealla's WebAssembly layer
    out_value(PrologArgs, Out),
    
    % 5. Construct the executable Goal
    Goal =.. [Verb | PrologArgs],
    
    % 6. ISO-Compliant Safe Call
    % Eliminates current_predicate/predicate_property bugs. If the command 
    % doesn't exist, catch/3 turns the engine error into a clean failure.
    catch(call(Goal), _, fail),
    !.

% Any commands that fail or do not exist return invalid
process_command(_, Out) :- Out = result(invalid).

% =========================================================================
% Helpers
% =========================================================================
% Helper: If an atom starts with an uppercase letter or is an underscore, treat it as a variable
bind_variables(Atom, _) :-
    atom(Atom),
    (   sub_atom(Atom, 0, 1, _, FirstChar),
        char_code(FirstChar, Code),
        Code >= 65, Code =< 90 % Matches 'A' through 'Z'
    ;   Atom == '_'
    ),
    !. 
bind_variables(Value, Value). 

% Optimized, independent filler remover
remove_fillers([], []).
remove_fillers([X|Xs], Ys) :-
    is_filler(X), !,
    remove_fillers(Xs, Ys).
remove_fillers([X|Xs], [X|Ys]) :-
    remove_fillers(Xs, Ys).

% Explicit filler words mapped directly as facts
is_filler(to).
is_filler(from).
is_filler(at).
is_filler(with).
is_filler(the).
is_filler(a).

