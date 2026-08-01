seed(4,4, cell_state(gate_or)). 

rows(5).
cols(5).

cell_state(out_of_bounds).
cell_state(super_position).
cell_state(empty).
cell_state(gate_and).
cell_state(gate_or).
cell_state(current_vertical).
cell_state(current_horizontal).


p(out_of_bounds, x).
p(super_position, s).
p(empty, e).
p(gate_and, a).
p(gate_or, o).
p(current_vertical, v).
p(current_horizontal, h).

% mod X, // Y
setup(TotalCells, RowCount, Grid) :-
    grid_rows(TotalCells, RowCount, TotalCells, FreshGrid),
    seed(Col, Row, CellState),
    update_grid(Row, Col, FreshGrid, CellState, Grid).


new_head(NewVal, [_| Rest], [NewVal | Rest]).

grid_rows(_, _, 0, []).
grid_rows(TotalCells, RowCount, Remaining, [Row | Rest]) :-
    Remaining > 0,
    Cols is TotalCells // RowCount,
    NuRemaining is Remaining - Cols,  
    row(Cols, Row),
    grid_rows(TotalCells, RowCount, NuRemaining, Rest), !.



row(0, []).
row(CellCounter, [cell_state(super_position) | Rest]) :-
    CellCounter > 0,
    NextCellCounter is CellCounter - 1,
    row(NextCellCounter, Rest).



% HELPERS

% Base Case: Index is 0. Swap the old Head out for the NewElement.
replace_nth(0, [_|Tail], NewElement, [NewElement|Tail]).

% Recursive Case: Keep the current Head, decrement Index, and move into the Tail.
replace_nth(Index, [Head|Tail], NewElement, [Head|NewTail]) :-
    Index > 0,
    NextIndex is Index - 1,
    replace_nth(NextIndex, Tail, NewElement, NewTail).

% Arguments: update_grid(RowIndex, ColIndex, Grid, NewValue, NewGrid)
update_grid(RowIdx, ColIdx, Grid, NewValue, NewGrid) :-
    % 1. Extract the target row
    nth0(RowIdx, Grid, OldRow),
    % 2. Create the updated row
    replace_nth(ColIdx, OldRow, NewValue, NewRow),
    % 3. Replace the old row with the new row in the main grid
    replace_nth(RowIdx, Grid, NewRow, NewGrid).

% 1. Main Predicate: Loops through each row
print_matrix([]). % Base case: nothing left to print
print_matrix([Row|Rest]) :-
    print_row(Row),
    nl,               % Move to the next line after finishing a row
    print_matrix(Rest).

% 2. Helper Predicate: Loops through elements in a single row
print_row([]).    % Base case: row is empty
print_row([cell_state(X)|Xs]) :-
    p(X, V),
    format(' ~w', [V]), % Prints element followed by a tab (~w = write, \t = tab)
    print_row(Xs).

% noop(Any) :- Any = [].
noop(_) :- nl, write('NOOP False: '), false.


% Example
% grid_rows(100, 10, 100, Grid), print_matrix(Grid), noop(Grid).
% setup(100, 10, Grid), print_matrix(Grid), noop(Grid).

% grid_rows(TotalCells, Rows, Remaining, Counter) :-
%     Remaining > 0,
%     Consume is TotalCells // Rows,
%     NuRemaining is Remaining - Consume,  
%     grid_rows(TotalCells, Rows, NuRemaining, RestCounter),
%     Counter is RestCounter + 1.
%
%
% grid_rows(_, _, 0, 0).



