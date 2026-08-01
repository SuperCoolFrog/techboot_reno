seed(2,2, gate_or). 

rows(5).
cols(5).

cell_state(out_of_bounds).
cell_state(super_position).
cell_state(empty).
cell_state(gate_and).
cell_state(gate_or).
cell_state(current_vertical).
cell_state(current_horizontal).

% mod X, // Y
setup(TotalCells, Rows, Grid) :-
    grid_rows(TotalCells, Rows, TotalCells, Grid).


grid_rows(TotalCells, Rows, Remaining, Counter) :-
    Remaining > 0,
    Consume is TotalCells // Rows,
    NuRemaining is Remaining - Consume,  
    grid_rows(TotalCells, Rows, NuRemaining, RestCounter),
    Counter is RestCounter + 1.


grid_rows(_, _, 0, 0).


% grid_rows(TotalCells, Rows, Remaining, Counter) :-
%     Remaining > 0,
%     Consume is TotalCells // Rows,
%     NuRemaining is Remaining - Consume,  
%     grid_rows(TotalCells, Rows, NuRemaining, RestCounter),
%     Counter is RestCounter + 1.
%
%
% grid_rows(_, _, 0, 0).



