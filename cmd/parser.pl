:- debug.
% Data

% .State Scene 3
state(3).

file_details(file_id(1), file_type(email), file_name(roy_1)).
file_details(file_id(2), file_type(email), file_name(roy_2)).
file_details(file_id(3), file_type(email), file_name(roy_3)).


fs(state(3),file_id(1)). 
fs(state(3),file_id(2)). 
fs(state(3),file_id(3)). 

shape(s, square).

connect(state(3), [missing(4, s)], [s,s,s,s,s,s]).


% Predicates

search(StateId, FT, FN) :- 
    state(StateId), 
    fs(state(StateId), file_id(FID)), 
    file_details(file_id(FID), file_type(FT), file_name(FN)).


missing_shape(missing(Idx, ShapeKey), Out) :-
    shape(ShapeKey, Shape),
    Out = result(idx, Shage).



connect_init(StateId, Missing, Puzzle) :- 
    state(StateId), 
    % connect(state(StateId), missing(Idx, ShapeKey), PuzzleList),
    connect(state(StateId), MissingList, PuzzleList),
    %shape(ShapeKey, Shape),
    Missing = MissingList, %result(missing, Idx, Shape),
    Puzzle = result(puzzle, 0, PuzzleList). 

% result(_) :- true.
%Missing is missing_shape(Idx, Shape). 

% findall(result(X, Y, Z), search(X, Y, Z), Result).
%
% ?- findall(result(X, Y, Z), search(X, Y, Z), Result).
% Result = [result(3, email, roy_1), result(3, email, roy_2), result(3, email, roy_3)].
% 
% ?- findall(result(X, Y, Z), search(X, Y, roy_1), Result).
% Result = [result(3, email, _)].
% 
% ?- findall(result(X, Y, Z), search(X, email, Z), Result).
% Result = [result(3, _, roy_1), result(3, _, roy_2), result(3, _, roy_3)].


