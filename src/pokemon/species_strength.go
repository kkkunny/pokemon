package pokemon

type SpeciesStrength [6]uint8

func (s SpeciesStrength) Sum() (sum uint16) {
	for _, v := range s {
		sum += uint16(v)
	}
	return sum
}

func (s SpeciesStrength) HP() uint8             { return s[0] }
func (s SpeciesStrength) Attack() uint8         { return s[1] }
func (s SpeciesStrength) Defense() uint8        { return s[2] }
func (s SpeciesStrength) SpecialAttack() uint8  { return s[3] }
func (s SpeciesStrength) SpecialDefense() uint8 { return s[4] }
func (s SpeciesStrength) Speed() uint8          { return s[5] }

func (s *SpeciesStrength) SetHP(v uint8)             { (*s)[0] = v }
func (s *SpeciesStrength) SetAttack(v uint8)         { (*s)[1] = v }
func (s *SpeciesStrength) SetDefense(v uint8)        { (*s)[2] = v }
func (s *SpeciesStrength) SetSpecialAttack(v uint8)  { (*s)[3] = v }
func (s *SpeciesStrength) SetSpecialDefense(v uint8) { (*s)[4] = v }
func (s *SpeciesStrength) SetSpeed(v uint8)          { (*s)[5] = v }
