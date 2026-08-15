puzzle_type(0, puzzle_intro).
puzzle_type(1, puzzle_easy).
puzzle_type(2, puzzle_med).
puzzle_type(3, puzzle_hard).

% Matches puzzle.go
gate_type(0, gate_empty).
gate_type(1, gate_unknown).
gate_type(2, gate_join).
gate_type(3, gate_split).
gate_type(4, gate_pass).

marker_type(0, marker_none).
marker_type(1, marker_yes).
marker_type(2, marker_no).


puzzle(config_id(1), state(17), puzzle_intro, gate_count(1), marker_count(1)).

% Helper
puzzle(Id) :- puzzle(Id, _,_,_,_).

% *************** Gates ******************
% OutputCols(26) ; OutputRows(36)
%
% X = OutputCols(26)/2 ; Y = OutputRows(36)/2

% ** Gate Valid
gate_valid(config_id(1), gate_idx(0), x(13), y(18), gate_pass). 

%% Helper
gate_valid(Id, [A,B,C,D]) :- gate_valid(Id, A,B,C,D).


% ** Gate Puzzle
gate_puzzle(config_id(1), gate_idx(0), x(13), y(18), gate_unknown). 

%% Helper
gate_puzzle(Id, [A,B,C,D]) :- gate_valid(Id, A,B,C,D).


% ** Gate Attempt
gate_attempt(config_id(1), gate_idx(0),  x(13), y(18), gate_unknown). 

%% Helper
gate_attempt(Id, [A,B,C,D]) :- gate_valid(Id, A,B,C,D).



% *************** Paths ******************

% ** Puzzle Paths
path_puzzle(config_id(1), start_x(13), start_y(35), end_x(13), end_y(15)).

%% Helper
path_puzzle(Id, [A, B, C, D]) :- path_puzzle(Id, A, B, C, D).

% ** Gate Paths
path_gate(gate_idx(0), gate_pass, start_x(13), start_y(16), end_x(13), end_y(0)).

%% Helper
path_gate(Id, [A, B, C, D, E]) :- path_puzzle(Id, A, B, C, D, E).


config_data(config(Id, State, PuzzleTypeId, GateCount, MarkerCount, PuzzlePaths, ValidGates, PuzzleGates, AttemptGates, GatePaths)) :-
    puzzle(config_id(Id), state(State), PuzzleType, gate_count(GateCount), marker_count(MarkerCount)),
    puzzle_type(PuzzleTypeId, PuzzleType),

    % Puzzle Paths
    findall(
        [StartX, StartY, EndX, EndY],
        path_puzzle(config_id(Id), start_x(StartX), start_y(StartY), end_x(EndX), end_y(EndY)),
        PuzzlePaths),

    % Valid Gates
    findall(
        [Idx, X, Y, GateTypeId],
        (gate_valid(config_id(Id), gate_idx(Idx), x(X), y(Y),  GateType), gate_type(GateTypeId, GateType)),
        ValidGates),
    % Puzzle Gates
    findall(
        [Idx, X, Y, GateTypeId],
        (gate_puzzle(config_id(Id), gate_idx(Idx), x(X), y(Y),  GateType), gate_type(GateTypeId, GateType)),
        PuzzleGates),
    % Attempt Gates
    findall(
        [Idx, X, Y, GateTypeId],
        (gate_attempt(config_id(Id), gate_idx(Idx), x(X), y(Y),  GateType), gate_type(GateTypeId, GateType)),
        AttemptGates),
    
    % Gate paths
    findall(
        [Idx, GateTypeId, StartX, StartY, EndX, EndY],
        (path_gate(gate_idx(Idx), GateType, start_x(StartX), start_y(StartY), end_x(EndX), end_y(EndY)), gate_type(GateTypeId, GateType)),
        GatePaths).

all_configs(Configs) :-
    findall(C, config_data(C), Configs).
