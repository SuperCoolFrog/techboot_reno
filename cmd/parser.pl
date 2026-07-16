% parse_command(+InputStr, -Action, -Arg)
parse_command(InputStr, Action, Arg) :-
    % 1. CRITICAL CONVERSION: Cast the bound Go string object cleanly to an Atom
    atom_string(InputAtom, InputStr),

    % 2. Remove the trailing '=' sign from the text stream
    atom_concat(CleanedCommand, '=', InputAtom),
    
    % 3. Split by space into Action and Arg atoms
    atomic_list_concat([ActionAtom, ArgAtom], ' ', CleanedCommand),
    
    % 4. Convert Action to lowercase (e.g., 'connect') to allow easy Go switches
    downcase_atom(ActionAtom, Action),
    
    % 5. Safely cast the text value into an actual native integer number
    atom_number(ArgAtom, Arg).

