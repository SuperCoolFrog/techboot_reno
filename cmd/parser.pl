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


% Predicates

search(FT, FN, StateId) :- 
    state(StateId), 
    fs(state(StateId), file_id(FID)), 
    file_details(file_id(FID), file_type(FT), file_name(FN)).

connect(Name, StateId) :-
		state(StateId),
		connection(state(StateId), Name).

% Parser
process_command(List) :-
    % 1. Clean filler words
    remove_fillers(List, CleanedList),
    
    % 2. Extract Verb and raw text arguments
    CleanedList = [Verb | RawArgs],
    
    % 3. Convert text placeholders (like 'FT') into true Prolog variables
    maplist(bind_variables, RawArgs, PrologArgs),
    
    % 4. Construct the executable Goal
    Goal =.. [Verb | PrologArgs],
    
    % 5. Execute the Goal safely
    current_predicate(_, Goal), 
    call(Goal),
    
    % 6. Print the answers to the user console
    print_results(RawArgs, PrologArgs).

% Helper: If an atom starts with an uppercase letter or is an underscore, treat it as a variable
bind_variables(Atom, Variable) :-
    atom(Atom),
    (   sub_atom(Atom, 0, 1, _, FirstChar),
        char_type(FirstChar, upper)
    ;   Atom == '_'
    ),
    !. % If it matches the criteria, leave Variable unbound.
bind_variables(Value, Value). % Otherwise, keep it as a concrete value (like 'my_file.txt').



% Parser
%process_command(List) :-
%    % Clean the list by removing filler words
%    remove_fillers(List, CleanedList),
%    
%    % Reconstruct into a Prolog Goal using the Univ (=..) operator
%    % CleanedList must look like [Verb, Arg1, Arg2, ...]
%    Goal =.. CleanedList,
%
%		current_predicate(_, Goal), call(Goal).
%    
%    % Check if the predicate exists and is true
%		% (   current_predicate(_, Goal), call(Goal)
%    % ->  format('Success: Fact "~w" is true.~n', [Goal])
%    % ;   format('Failure: Fact "~w" is false or undefined.~n', [Goal])
%    % ).
%
%% 3. Helper: Strip out structural filler words dynamically
%remove_fillers([], []).
%remove_fillers([X|Xs], Ys) :-
%    member(X, [to, from, at, with, the, a]), !, % Add any filler words here
%    remove_fillers(Xs, Ys).
%remove_fillers([X|Xs], [X|Ys]) :-
%    remove_fillers(Xs, Ys).
%
%
