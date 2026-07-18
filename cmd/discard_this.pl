% 1. Sample Database Facts (Any predicate structure works)
connect(3, 43).
disconnect(5, 99).
move(player1, room2).

% 2. Main Entry Point
process_command(List) :-
    % Clean the list by removing filler words
    remove_fillers(List, CleanedList),
    
    % Reconstruct into a Prolog Goal using the Univ (=..) operator
    % CleanedList must look like [Verb, Arg1, Arg2, ...]
    Goal =.. CleanedList,
    
    % Check if the predicate exists and is true
    (   current_predicate(_, Goal), call(Goal)
    ->  format('Success: Fact "~w" is true.~n', [Goal])
    ;   format('Failure: Fact "~w" is false or undefined.~n', [Goal])
    ).

% 3. Helper: Strip out structural filler words dynamically
remove_fillers([], []).
remove_fillers([X|Xs], Ys) :-
    member(X, [to, from, at, with, the, a]), !, % Add any filler words here
    remove_fillers(Xs, Ys).
remove_fillers([X|Xs], [X|Ys]) :-
    remove_fillers(Xs, Ys).
